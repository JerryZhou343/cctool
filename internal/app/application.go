package app

import (
	"context"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/merge"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/store/aliyun"
	"github.com/JerryZhou343/cctool/internal/translate/baidu"
	"github.com/JerryZhou343/cctool/internal/translate/google"
	"github.com/JerryZhou343/cctool/internal/translate/tencent"
	"os"
	"strings"
	"sync"

	aliSpeech "github.com/JerryZhou343/cctool/internal/text/aliyun"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
	"github.com/pkg/errors"
	"path/filepath"
	"strconv"
	"time"
)

type Application struct {
	//所有的翻译工具
	translatorSet map[string]*Translator
	//空闲中的工具
	idleTranslator map[string]struct{}
	//待清理工具
	cleanTranslator map[string]struct{}
	translatorLock  *sync.Mutex

	translateTaskChan chan *TranslateTask
	ctx               context.Context
	cancelFunc        context.CancelFunc
	msgChan           chan string
}

func NewApplication() *Application {
	ret := &Application{
		translatorSet:     map[string]*Translator{},
		idleTranslator: map[string]struct{}{},
		cleanTranslator:map[string]struct{}{},
		translateTaskChan: make(chan *TranslateTask, 1000),
		translatorLock:    new(sync.Mutex),

		msgChan:           make(chan string, 1000),
	}
	ret.ctx, ret.cancelFunc = context.WithCancel(context.Background())
	return ret
}

func (a *Application) Destroy() {
	a.cancelFunc()
}

func (a *Application) Run() {
	go a.translate()
}

func (a *Application) GetRunningMsg() string {
	msg := <-a.msgChan
	return msg
}

func (a *Application) LoadTranslateTools() (err error) {
	err = conf.Load()
	if err != nil {
		return
	}

	a.translatorLock.Lock()
	defer a.translatorLock.Unlock()

	for _, itr := range conf.G_Config.TransTools {
		a.cleanTranslator[itr] = struct{}{}
	}

	for _, itr := range conf.G_Config.TransTools {
		if _, ok := a.translatorSet[itr]; !ok {
			switch itr {
			case "google":
				a.translatorSet[itr] = NewTranslator(itr, google.NewTranslator(),
					time.Duration(conf.G_Config.Google.Interval)*time.Millisecond)
			case "baidu":
				if conf.G_Config.Baidu.Check() {
					a.translatorSet[itr] = NewTranslator(itr,
						baidu.NewTranslator(conf.G_Config.Baidu.AppId, conf.G_Config.Baidu.SecretKey),
						time.Duration(conf.G_Config.Baidu.Interval)*time.Millisecond)
				}
			case "tencent":
				if conf.G_Config.Tencent.Check() {
					a.translatorSet[itr] = NewTranslator(itr,
						tencent.NewTranslator(conf.G_Config.Tencent.Qtk, conf.G_Config.Tencent.Qtv),
						time.Duration(conf.G_Config.Tencent.Interval)*time.Millisecond)
				}
			}
			a.idleTranslator[itr] = struct{}{}
		}
		// 不需要清理
		if _, ok := a.cleanTranslator[itr]; ok {
			delete(a.cleanTranslator, itr)
		}
	}

	return nil

}

func (a *Application) AddTranslateTask(task *TranslateTask) (err error) {
	a.translateTaskChan <- task
	a.msgChan <- fmt.Sprintf("添加任务成功 %s",task)
	return nil
}

func (a *Application) translate() {
	for {
		select {
		case task := <-a.translateTaskChan:
			func() {
				for {
					a.translatorLock.Lock()
					//寻找一个空闲的translator 执行翻译任务
					if len(a.idleTranslator) > 0 {
						for k, _ := range a.idleTranslator {
							if v, ok := a.translatorSet[k]; ok {
								if !v.Running {
									v.Start()
									go v.Do(a.ctx, task, a.msgChan, a.translateTaskDone)
									delete(a.idleTranslator,k)
								}

								break
							}
						}
						a.translatorLock.Unlock()
						break
					} else {
						a.translatorLock.Unlock()
						time.Sleep(1 * time.Second)
					}

				}
			}()

		case <-a.ctx.Done():
			return
		}
	}
}

//翻译任务做完后回调
func (a *Application) translateTaskDone(translator *Translator) {
	a.translatorLock.Lock()
	defer a.translatorLock.Unlock()
	translator.Done()
	if _, ok := a.cleanTranslator[translator.Name]; ok {
		delete(a.translatorSet, translator.Name)
		delete(a.cleanTranslator, translator.Name)
	} else {
		a.idleTranslator[translator.Name] = struct{}{}
	}

}

func (a *Application) Merge() error {
	engine := merge.NewMerge()
	return engine.Merge(flags.MergeStrategy, flags.DstFile, flags.SrcFiles...)
}

func (a *Application) GenerateSrt(video string, channelId int) (err error) {
	var (
		uri      string
		ret      []*srt.Srt
		absVideo string
		objName  string
	)
	//1. 抽取音频
	extractor := voice.NewExtractor(strconv.Itoa(conf.G_Config.SampleRate))
	err = extractor.Valid()
	if err != nil {
		return
	}
	absVideo, err = filepath.Abs(video)
	if err != nil {
		return errors.Wrapf(status.ErrReadFileFailed, "%s", video)
	}
	flag := utils.CheckFileExist(absVideo)
	if !flag {
		return errors.Wrapf(status.ErrFileNotExits, "%s", video)
	}
	fileName := filepath.Base(absVideo)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	flag = utils.CheckFileExist(conf.G_Config.AudioCachePath)
	if !flag {
		err = os.MkdirAll(conf.G_Config.AudioCachePath, os.ModePerm)
		if err != nil {
			return
		}
	}
	dstAudioFile := filepath.Join(conf.G_Config.AudioCachePath, name+".mp3")
	err = extractor.ExtractAudio(absVideo, dstAudioFile)
	if err != nil {
		return
	}
	//2. 存储
	storage := aliyun.NewAliyunOSS(conf.G_Config.Aliyun.OssEndpoint,
		conf.G_Config.Aliyun.AccessKeyId, conf.G_Config.Aliyun.AccessKeySecret,
		conf.G_Config.Aliyun.BucketName, conf.G_Config.Aliyun.BucketDomain)
	uri, objName, err = storage.UploadFile(dstAudioFile)
	//3. 识别
	speech := aliSpeech.NewSpeech(conf.G_Config.Aliyun.AccessKeyId, conf.G_Config.Aliyun.AccessKeySecret,
		conf.G_Config.Aliyun.AppKey)
	ret, err = speech.Recognize(uri, channelId)
	if err != nil {
		return
	}
	//4. 输出
	srtDstFilePath := filepath.Join(conf.G_Config.SrtPath, name+".srt")
	srt.WriteSrt(srtDstFilePath, ret)

	//清理文件
	os.Remove(dstAudioFile)
	storage.DeleteFile(objName)
	return nil
}
