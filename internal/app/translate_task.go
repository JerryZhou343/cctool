package app

import (
	"fmt"
	"path/filepath"
	"strings"
)

type TranslateTask struct {
	SrcFile    string
	DstFile    string
	From       string
	To         string
	Progress   float32
	Merge      bool // 双语字幕
	State      TaskState
	translator string
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
	return fmt.Sprintf("[翻译任务] 工具: %s 源文件: %s 目标文件: %s 进度: %.2f", t.translator, t.SrcFile, t.DstFile, t.Progress*100) +
		"%100 " + fmt.Sprintf("状态: %s", t.State)
}

func (t *TranslateTask) Init() (err error) {
	var (
		absFilePath string
	)
	absFilePath, err = filepath.Abs(t.SrcFile)
	if err != nil {
		return
	}
	absPath := filepath.Dir(absFilePath)
	fileName := filepath.Base(absFilePath)
	ext := filepath.Ext(fileName)
	name := strings.Trim(fileName, ext)
	t.DstFile = filepath.Join(absPath, fmt.Sprintf("%s_%s%s", name, t.To, ext))
	return
}
