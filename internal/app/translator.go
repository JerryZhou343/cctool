package app

import (
	"context"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/translate"
	"path/filepath"
	"strings"
	"time"
)

type Translator struct {
	interval time.Duration //请求间隔
	tool     translate.Translate
	Name     string
	Running  bool
}

func NewTranslator(name string, tool translate.Translate, interval time.Duration) *Translator {
	return &Translator{
		interval: interval,
		tool:     tool,
		Name:     name,
	}
}

func (t *Translator) Start() {
	t.Running = true
}

func (t *Translator) Done() {
	t.Running = false
}

func (t *Translator) Do(ctx context.Context, task *TranslateTask, msg chan string, doneCallBack func(*Translator)) {
	var (
		err         error
		src         []*srt.Srt
		absFilePath string
		subtitle    string
		subtitleSet map[int]string
	)
	defer doneCallBack(t)
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
		tryTimes := 10
		select {
		default:
			for {
				if tryTimes == 0{
					return
				}
				time.Sleep(t.interval)
				subtitle, err = t.tool.Do(itr.Subtitle, task.From, task.To)
				if err != nil {
					task.State = TaskStateFailed
					tryTimes--
					msg <- fmt.Sprintf("translator: %s,task: %s err:%+v",t.Name, task, err)
				} else {
					msg <- fmt.Sprintf("translator: %s,task: %s ",t.Name, task)
					subtitleSet[itr.Sequence] = subtitle
					task.Progress = float32(idx) / float32(total)
					break
				}
			}

		case <-ctx.Done():
			return
		}

	}

	//3.目标文件输出
	//3.1 计算目标文件名
	absFilePath, err = filepath.Abs(task.SrcFile)
	if err != nil {
		msg <- fmt.Sprintf("translator: %s,task: %s err:%+v",t.Name, task, err)
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

	err = srt.WriteSrt(dstFile, src)
	if err != nil {
		msg <- fmt.Sprintf("translator: %s,task: %s err:%+v",t.Name, task, err)
		return
	}

	task.State = TaskStateDone
	msg <- fmt.Sprintf("translator: %s,task: %s ",t.Name, task)
}
