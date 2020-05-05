package app

import (
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
	"github.com/JerryZhou343/cctool/internal/translate/baidu"
	"github.com/JerryZhou343/cctool/internal/translate/google"
	"github.com/JerryZhou343/cctool/internal/translate/tencent"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
	"github.com/pkg/errors"
	"path/filepath"
	"strconv"
	"time"
)

type Application struct {
	translator translate.Translate
	interval   time.Duration
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Translate() (err error) {
	var (
		ret []*srt.Srt
	)
	switch flags.TransTool {
	case flags.TransTool_Baidu:
		a.translator = baidu.NewTranslator(&conf.G_Config.Baidu)
		a.interval = time.Millisecond * time.Duration(conf.G_Config.Baidu.Interval)
	case flags.TransTool_Google:
		a.translator = google.NewTranslator()
		a.interval = time.Millisecond * time.Duration(conf.G_Config.Google.Interval)
	case flags.TransTool_Tencent:
		a.translator = tencent.NewTranslator(conf.G_Config.Tencent.Qtk, conf.G_Config.Tencent.Qtv)
		a.interval = time.Millisecond * time.Duration(conf.G_Config.Tencent.Interval)
	default:
		return status.ErrInitTranslatorFailed
	}

	for _, itr := range flags.SrcFiles {
		ret, err = srt.Open(itr)
		if err != nil {
			return err
		}
		err = a.translate(itr, ret)
		if err != nil {
			return
		}
	}

	return
}

func (a *Application) translate(srcPath string, src []*srt.Srt) (err error) {
	var (
		tmpResult   string
		absFilePath string
		ret         map[int]string
	)
	absFilePath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	absPath := filepath.Dir(absFilePath)
	fileName := filepath.Base(absFilePath)
	dstFile := filepath.Join(absPath, flags.To+"_"+fileName)

	ret = make(map[int]string)
	for _, itr := range src {
		tmpResult, err = a.translator.Do(itr.Subtitle, flags.From, flags.To)
		if err != nil {
			return
		}
		ret[itr.Sequence] = tmpResult
		time.Sleep(a.interval)
	}

	if flags.Merge {
		for _, itr := range src {
			if v, ok := ret[itr.Sequence]; ok {
				itr.Subtitle = v + "\n" + itr.Subtitle
			}
		}
		srt.WriteSrt(dstFile, src)
	} else {
		for _, itr := range src {
			if v, ok := ret[itr.Sequence]; ok {
				itr.Subtitle = v
			}
		}
		srt.WriteSrt(dstFile, src)
	}
	return
}

func (a *Application) Merge() error {
	engine := merge.NewMerge()
	return engine.Merge(flags.MergeStrategy, flags.DstFile, flags.SrcFiles...)
}

func (a *Application) GenerateSrt(video string, channelId int) (err error) {
	var (
		uri string
		ret []*srt.Srt
		absVideo string
		objName string
	)
	//1. 抽取音频
	extractor := voice.NewExtractor(strconv.Itoa(conf.G_Config.SampleRate))
	err = extractor.Valid()
	if err != nil {
		return
	}
	absVideo, err  = filepath.Abs(video)
	if err != nil{
		return errors.Wrapf( status.ErrReadFileFailed,"%s",video)
	}
	flag := utils.CheckFileExist(absVideo)
	if !flag {
		return errors.Wrapf(status.ErrFileNotExits, "%s", video)
	}
	fileName := filepath.Base(absVideo)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	dstAudioFile := filepath.Join(conf.G_Config.AudioCachePath , name + ".mp3")
	extractor.ExtractAudio(absVideo, dstAudioFile)
	//2. 存储
	storage := aliyun.NewAliyunOSS(conf.G_Config.Aliyun.OssEndpoint,
		conf.G_Config.Aliyun.AccessKeyId, conf.G_Config.Aliyun.AccessKeySecret,
		conf.G_Config.Aliyun.BucketName, conf.G_Config.Aliyun.BucketDomain)
	uri,objName, err = storage.UploadFile(dstAudioFile)
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
