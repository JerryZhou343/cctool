package app

import (
	"context"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/translate"
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

func (t *Translator) Do(ctx context.Context, task *TranslateTask, doneCallBack func(*Translator)) {
	var (
		err         error
		src         []*srt.Srt
		subtitle    string
		subtitleSet map[int]string
	)
	defer doneCallBack(t)
	//1.准备数据
	err = task.Init()
	if err != nil {
		task.Failed(err)
		return
	}

	subtitleSet = make(map[int]string)
	src, err = srt.Open(task.SrcFile)
	if err != nil {

		task.Failed(err)
		return
	}

	total := len(src)
	task.State = TaskStateDoing
	task.translator = t.Name
	//2.翻译
	for idx, itr := range src {
		tryTimes := 10
		select {
		default:
			for {
				if tryTimes == 0 {
					task.State = TaskStateFailed
					return
				}
				time.Sleep(t.interval)
				tmp := strings.ReplaceAll(strings.ReplaceAll(itr.Subtitle, "\r\n", " "), "\n", " ")
				strings.TrimSpace(tmp)
				if tmp != ""{
					subtitle, err = t.tool.Do(tmp, task.From, task.To)
				}
				if err != nil {
					task.Failed(err)

					task.State = TaskStateTrying
					tryTimes--
				} else {
					subtitleSet[itr.Sequence] = subtitle
					task.Progress = float32(idx+1) / float32(total)
					task.State = TaskStateDoing
					break
				}
			}

		case <-ctx.Done():
			return
		}

	}

	//3.目标文件输出
	for _, itr := range src {
		if v, ok := subtitleSet[itr.Sequence]; ok {
			if task.Merge {
				itr.Subtitle = v + "\r\n" + itr.Subtitle
			} else {
				itr.Subtitle = v
			}
		}
	}

	err = srt.WriteSrt(task.DstFile, src)
	if err != nil {
		task.Failed(err)
		return
	}
	task.State = TaskStateDone
	task.Failed(nil)
}
