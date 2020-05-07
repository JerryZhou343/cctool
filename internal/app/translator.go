package app

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/translate"
	"path/filepath"
	"strings"
	"time"
)

type TranslateTask struct {
	SrcFile  string
	DstFile  string
	From     string
	To       string
	Progress float32
	Merge    bool // 双语字幕
	State    TaskState
}

type Translator struct {
	interval time.Duration //请求间隔
	tool     translate.Translate
	Name     string
}

func (t *Translator) Do(task *TranslateTask, msg chan string) {
	var (
		err         error
		src         []*srt.Srt
		absFilePath string
		subtitle    string
		subtitleSet map[int]string
	)

	//1.准备数据
	subtitleSet = make(map[int]string)
	src, err = srt.Open(task.SrcFile)
	if err != nil {
		return
	}
	total := len(src)
	task.State = TaskStateDoing
	//2.翻译
	for idx, itr := range src {
		subtitle, err = t.tool.Do(itr.Subtitle, task.From, task.To)
		if err != nil {
			task.State = TaskStateFailed
			return
		}
		subtitleSet[itr.Sequence] = subtitle
		task.Progress = float32(idx) / float32(total)
	}

	//3.目标文件输出
	//3.1 计算目标文件名
	absFilePath, err = filepath.Abs(task.SrcFile)
	if err != nil {
		return
	}
	absPath := filepath.Dir(absFilePath)
	fileName := filepath.Base(absFilePath)
	ext := filepath.Ext(fileName)
	name := strings.Trim(fileName, ext)
	dstFile := filepath.Join(absPath, fmt.Sprintf("%s_%s.%s", name, task.To, ext))

	//3.2 内容合并
	for _, itr := range src {
		if v, ok := subtitleSet[itr.Sequence]; ok {
			if task.Merge {
				itr.Subtitle = v + "\r\n" + itr.Subtitle
			} else {
				itr.Subtitle = v
			}
		}
	}
	srt.WriteSrt(dstFile, src)
	task.State = TaskStateDone
}
