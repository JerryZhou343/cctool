package app

import (
	"context"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/store"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SrtGenerator struct {
	Name    string
	storage store.Store
	speech  text.ISpeech

	Running bool
}

func NewSrtGenerator(name string, storage store.Store, speech text.ISpeech) *SrtGenerator {
	return &SrtGenerator{
		Name:    name,
		storage: storage,
		speech:  speech,
	}
}

func (s *SrtGenerator) Start() {
	s.Running = true
}

func (s *SrtGenerator) Done() {
	s.Running = false
}

func (s *SrtGenerator) Do(ctx context.Context, task *GenerateTask, doneCallBack func(generator *SrtGenerator)) {
	var (
		uri      string
		ret      []*srt.Srt
		absVideo string
		objName  string
		err      error
	)
	defer doneCallBack(s)
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
	extractor := voice.NewExtractor(strconv.Itoa(conf.G_Config.SampleRate), conf.G_Config.FFmpeg)
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
	uri, objName, err = s.storage.UploadFile(dstAudioFile)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		return
	}
	defer s.storage.DeleteFile(objName)

	//3. 识别
	task.Step = GenerateStepRecognize
	ret, err = s.speech.Recognize(ctx, uri, task.ChannelId)
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
}
