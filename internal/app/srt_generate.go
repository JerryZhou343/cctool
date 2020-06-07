package app

import (
	"context"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/store"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/JerryZhou343/cctool/internal/voice"
	"github.com/sirupsen/logrus"
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
		wret     []*srt.Srt // 词断句结果
		sret     []*srt.Srt //原始句子结果
		absVideo string
		objName  string
		err      error
	)
	defer doneCallBack(s)
	logrus.Infof("get generate task %s",task)
	//前置检查
	absVideo, err = filepath.Abs(task.SrcFile)
	if err != nil {
		task.State = TaskStateFailed
		task.Failed(err)
		logrus.Errorf("%s check file path failed [%+v]", task, err)
		return
	}
	flag := utils.CheckFileExist(absVideo)
	if !flag {

		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("%s check file exists failed [%+v]", task, err)
		return
	}

	fileName := filepath.Base(absVideo)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	dstAudioFile := filepath.Join(conf.G_Config.AudioCachePath, name+".mp3")
	wsrtDstFilePath := filepath.Join(conf.G_Config.SrtPath, name+"_word.srt")
	ssrtDstFilePath := filepath.Join(conf.G_Config.SrtPath, name+"_sentence.srt")
	task.DstFile = wsrtDstFilePath

	task.State = TaskStateDoing
	logrus.Infof("start extract audio [%s]",task)
	//1. 抽取音频
	task.Step = GenerateStepAudio
	extractor := voice.NewExtractor(strconv.Itoa(int(conf.G_Config.SampleRate)), conf.G_Config.FFmpeg)
	err = extractor.Valid()
	if err != nil {
		task.State = TaskStateFailed
		logrus.Errorf("task:[%s], check extractor  failed [%v]", task, err)
		return
	}
	flag = utils.CheckFileExist(conf.G_Config.AudioCachePath)
	if !flag {
		err = os.MkdirAll(conf.G_Config.AudioCachePath, os.ModePerm)
		if err != nil {
			task.Failed(err)
			task.State = TaskStateFailed
			logrus.Errorf("create directory failed [%v]", err)
			return
		}
	}

	flag = utils.CheckFileExist(dstAudioFile)
	if flag {
		err = os.Remove(dstAudioFile)
		if err != nil {
			logrus.Errorf("remove failed [%s] failed [%v]", dstAudioFile, err)
			task.Failed(err)
			task.State = TaskStateFailed
			return
		}
	}

	err = extractor.ExtractAudio(absVideo, dstAudioFile)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("task[%s] extract audio failed [%v]", task, err)
		return
	}
	defer os.Remove(dstAudioFile)
	logrus.Infof("end extract audio %s",task)

	//2. 存储
	task.Step = GenerateStepOss
	logrus.Infof("start upload file [%s]",task)
	uri, objName, err = s.storage.UploadFile(dstAudioFile)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("task[%s] upload file failed [%v]", task, err)
		return
	}
	defer s.storage.DeleteFile(objName)
	logrus.Infof("end upload file [%s]",task)
	//3. 识别
	task.Step = GenerateStepRecognize
	logrus.Infof("start recognize [%s]",task)
	sret, wret, err = s.speech.Recognize(ctx, uri)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("task[%s] recognize failed [%v]", task, err)
		return
	}
	logrus.Infof("end recognize [%s] result sret [%d], wret[%d]",task,len(sret), len(wret))
	//4. 输出
	task.Step = GenerateStepGenerateSrt
	err = srt.WriteSrt(wsrtDstFilePath, wret)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("task[%s] write srt failed [%v]", task, err)
		return
	}
	err = srt.WriteSrt(ssrtDstFilePath, sret)
	if err != nil {
		task.Failed(err)
		task.State = TaskStateFailed
		logrus.Errorf("task[%s] write srt failed [%v]", task, err)
		return
	}
	logrus.Infof("write srt end [%s]",task)
	task.State = TaskStateDone
}
