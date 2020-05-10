package app

import (
	"fmt"
	"path/filepath"
	"strings"
)

type TaskState int

const (
	TaskStateUnknown = iota
	TaskStateInit    = 1
	TaskStateDoing   = 2
	TaskStateDone    = 3
	TaskStateFailed  = 4
)

func (t TaskState) String() string {
	switch t {
	case TaskStateDoing:
		return "正在进行"
	case TaskStateDone:
		return "已完成"
	case TaskStateFailed:
		return "失败"
	default:
		return "未知状态"
	}
}

//字幕翻译任务
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
	return fmt.Sprintf("[翻译] 工具: %s 源文件: %s 目标文件: %s 进度: %.2f", t.translator, t.SrcFile, t.DstFile, t.Progress*100) +
		"% " + fmt.Sprintf("状态: %s", t.State)
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

//字幕格式转换任务
type ConvertTask struct {
	SrcFile string
	DstFile string
	State   TaskState
	From    string
	To      string
}

func NewConvertTask(srcFile, from, to string) *ConvertTask {
	return &ConvertTask{
		SrcFile: srcFile,
		DstFile: from,
		State:   TaskStateInit,
		From:    from,
		To:      to,
	}
}

func (c *ConvertTask) String() string {
	return fmt.Sprintf("[字幕格式转换] From: %s To: %s 源文件: %s 目标文件：%s  状态: %s",
		c.From, c.To, c.SrcFile, c.DstFile, c.State)
}

func (c *ConvertTask) Init() (err error) {
	var (
		absFilePath string
	)
	absFilePath, err = filepath.Abs(c.SrcFile)
	if err != nil {
		return
	}
	absPath := filepath.Dir(absFilePath)
	fileName := filepath.Base(absFilePath)
	ext := filepath.Ext(fileName)
	name := strings.Trim(fileName, ext)
	c.DstFile = filepath.Join(absPath, fmt.Sprintf("%s.%s", name, c.To))
	return
}
