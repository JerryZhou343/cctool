package app

import (
	"context"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/merge"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/store/aliyun"
	"os"
	"strings"

	aliSpeech "github.com/JerryZhou343/cctool/internal/text/aliyun"
	"github.com/JerryZhou343/cctool/internal/translate"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"path/filepath"
	"strconv"
	"time"
)

type Application struct {
	translatorSet     map[string]translate.Translate
	interval          time.Duration
	translateTaskChan chan *TranslateTask
	ctx               context.Context
}

func NewApplication() *Application {
	ret := &Application{
		translatorSet:     map[string]translate.Translate{},
		translateTaskChan: make(chan *TranslateTask, 1000),
	}

	return ret
}

func (a *Application) Destroy() {
	ants.Release()
}

func (a *Application) Run() {
	go a.translate()
}

func (a *Application) AddTranslateTask(task *TranslateTask) (err error) {
	a.translateTaskChan <- task
	return nil
}

func (a *Application) translate() {

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
