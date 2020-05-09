package app

import "fmt"

type TranslateTask struct {
	SrcFile  string
	DstFile  string
	From     string
	To       string
	Progress float32
	Merge    bool // 双语字幕
	State    TaskState
}

func NewTranslateTask(srcFile, from, to string, merge bool) *TranslateTask {
	return &TranslateTask{
		SrcFile:  srcFile,
		DstFile:  "",
		From:     from,
		To:       to,
		Progress: 0,
		Merge:    merge,
		State:    TaskStateInit,
	}
}

func (t *TranslateTask) String() string {
	return fmt.Sprintf("[翻译任务]源文件: %s, 目标文件: %s 进度: %.2f",t.SrcFile, t.DstFile, t.Progress*100) +
		"%100 "+ fmt.Sprintf("状态: %s",t.State)
}
