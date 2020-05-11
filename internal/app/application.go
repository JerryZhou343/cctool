package app

import (
	"context"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/bcc"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/convert"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/merge"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/store/aliyun"
	"github.com/JerryZhou343/cctool/internal/translate/baidu"
	"github.com/JerryZhou343/cctool/internal/translate/google"
	"github.com/JerryZhou343/cctool/internal/translate/tencent"
	"github.com/panjf2000/ants/v2"
	"os"
	"strings"
	"sync"

	aliSpeech "github.com/JerryZhou343/cctool/internal/text/aliyun"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
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
	//翻译任务
	translateTaskChan chan Task
	//转换任务
	convertTaskChan chan Task
	//生成字幕任务
	generateTaskChan chan Task

	//任务数组
	taskSlice []Task
	//
	ctx        context.Context
	cancelFunc context.CancelFunc
	msgChan    chan string
}

func NewApplication() *Application {
	ret := &Application{
		translatorSet:     map[string]*Translator{},
		idleTranslator:    map[string]struct{}{},
		cleanTranslator:   map[string]struct{}{},
		translateTaskChan: make(chan Task, 1000),
		translatorLock:    new(sync.Mutex),

		convertTaskChan: make(chan Task, 100),

		generateTaskChan: make(chan Task, 100),

		msgChan: make(chan string, 1000),
	}
	ret.ctx, ret.cancelFunc = context.WithCancel(context.Background())
	return ret
}

func (a *Application) Destroy() {
	a.cancelFunc()
	ants.Release()
}

func (a *Application) Run() {
	go a.translate()
	go a.convert()
	go a.generate()
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

func (a *Application) AddTask(task Task) (err error) {
	switch task.Type() {
	case TaskTypeTranslate:
		a.translateTaskChan <- task
	case TaskTypeGenerate:
		a.generateTaskChan <- task
	case TaskTypeConvert:
		a.convertTaskChan <- task
	}
	a.msgChan <- fmt.Sprintf("添加任务成功 %s", task)
	a.taskSlice = append(a.taskSlice, task)
	return nil
}

func (a *Application) CheckTask() {
	var (
		allDone bool
	)
	for {
		allDone = true
		select {
		case <-time.After(2 * time.Second):
			for _, itr := range a.taskSlice {
				a.msgChan <- fmt.Sprintf("%s", itr)

				//任务超过最大重试次数就不再尝试
				if itr.GetState() == TaskStateFailed && itr.GetFailedTimes() < 10 {
					a.AddTask(itr)
				}
				if itr.GetState() != TaskStateDone && itr.GetFailedTimes() < 10 {
					allDone = false
				}
			}

			if allDone {
				return
			}

		case <-a.ctx.Done():
			return
		}
	}
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
									go v.Do(a.ctx, task.(*TranslateTask), a.translateTaskDone)
									delete(a.idleTranslator, k)
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

func (a *Application) generate() {
	for {
		select {
		case t := <-a.generateTaskChan:
			a.generateSrt(t.(*GenerateTask))
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *Application) generateSrt(task *GenerateTask) {
	var (
		uri      string
		ret      []*srt.Srt
		absVideo string
		objName  string
		err      error
	)
	//前置检查
	absVideo, err = filepath.Abs(task.SrcFile)
	if err != nil {
		task.State = TaskStateFailed
		task.Failed(err)
		return
	}
	flag := utils.CheckFileExist(absVideo)
	if !flag {

		task.Failed(err)
		task.State = TaskStateFailed
		return
	}

	fileName := filepath.Base(absVideo)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	dstAudioFile := filepath.Join(conf.G_Config.AudioCachePath, name+".mp3")
	srtDstFilePath := filepath.Join(conf.G_Config.SrtPath, name+".srt")
	task.DstFile = srtDstFilePath

	task.State = TaskStateDoing
	//1. 抽取音频
	task.Step = GenerateStepAudio
	extractor := voice.NewExtractor(strconv.Itoa(conf.G_Config.SampleRate))
	err = extractor.Valid()
	if err != nil {
		task.State = TaskStateFailed
		return
	}
	flag = utils.CheckFileExist(conf.G_Config.AudioCachePath)
	if !flag {
		err = os.MkdirAll(conf.G_Config.AudioCachePath, os.ModePerm)
		if err != nil {

			task.Failed(err)
			task.State = TaskStateFailed
			return
		}
	}

	flag = utils.CheckFileExist(dstAudioFile)
	if flag {
		err = os.Remove(dstAudioFile)
		if err != nil {

			task.Failed(err)
			task.State = TaskStateFailed
			return
		}
	}

	err = extractor.ExtractAudio(absVideo, dstAudioFile)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		return
	}
	defer os.Remove(dstAudioFile)

	//2. 存储
	task.Step = GenerateStepOss
	storage := aliyun.NewAliyunOSS(conf.G_Config.Aliyun.OssEndpoint,
		conf.G_Config.Aliyun.AccessKeyId, conf.G_Config.Aliyun.AccessKeySecret,
		conf.G_Config.Aliyun.BucketName, conf.G_Config.Aliyun.BucketDomain)
	uri, objName, err = storage.UploadFile(dstAudioFile)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		return
	}
	defer storage.DeleteFile(objName)

	//3. 识别
	task.Step = GenerateStepRecognize
	speech := aliSpeech.NewSpeech(conf.G_Config.Aliyun.AccessKeyId, conf.G_Config.Aliyun.AccessKeySecret,
		conf.G_Config.Aliyun.AppKey)
	ret, err = speech.Recognize(uri, task.ChannelId)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		return
	}
	//4. 输出
	task.Step = GenerateStepGenerateSrt
	err = srt.WriteSrt(srtDstFilePath, ret)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
	}
	task.State = TaskStateDone
	return
}

//字幕格式转换
func (a *Application) convert() {
	for {
		select {
		case t := <-a.convertTaskChan:
			_ = ants.Submit(func() {
				var (
					err error
					src *bcc.BCC
					ret []*srt.Srt
				)
				task := t.(*ConvertTask)
				task.State = TaskStateInit
				err = task.Init()
				if err != nil {
					task.State = TaskStateFailed
					task.Failed(err)
					return
				}

				src, err = bcc.Open(task.SrcFile)
				if err != nil {
					task.State = TaskStateFailed
					task.Failed(err)
				}

				//doing
				task.State = TaskStateDoing
				ret = convert.BCC2SRT(src)

				//done
				err = srt.WriteSrt(task.DstFile, ret)
				if err != nil {
					task.State = TaskStateFailed
					task.Failed(err)
					return
				}
				task.State = TaskStateDone
			})
		case <-a.ctx.Done():
			return
		}
	}
}
