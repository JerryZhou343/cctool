package app

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Task interface {
	GetState() TaskState
	String() string
	Type() TaskType
	GetFailedTimes() int
	Failed(err error)
}

type TaskType int

const (
	TaskTypeUnknown = iota
	TaskTypeTranslate
	TaskTypeGenerate
	TaskTypeConvert
	TaskTypeMerge
	TaskTypeClean
)

type TaskState int

const (
	TaskStateUnknown = iota
	TaskStateInit
	TaskStateDoing
	TaskStateTrying
	TaskStateDone
	TaskStateFailed
)

func (t TaskState) String() string {
	switch t {
	case TaskStateInit:
		return "初始化"
	case TaskStateDoing:
		return "正在进行"
	case TaskStateDone:
		return "已完成"
	case TaskStateFailed:
		return "失败"
	case TaskStateTrying:
		return "正在重试"
	default:
		return "未知状态"
	}
}

//字幕翻译任务
type TranslateTask struct {
	SrcFile     string
	DstFile     string
	From        string
	To          string
	Progress    float32
	Merge       bool // 双语字幕
	State       TaskState
	translator  string
	FailedTimes int
	Err         error
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
	if t.Err == nil {
		return fmt.Sprintf("[翻译] 工具: %s 源文件: %s 目标文件: %s 进度: %.2f", t.translator, t.SrcFile, t.DstFile, t.Progress*100) +
			"% " + fmt.Sprintf("状态: %s", t.State)
	} else {
		return fmt.Sprintf("[翻译] 工具: %s 源文件: %s 目标文件: %s 进度: %.2f", t.translator, t.SrcFile, t.DstFile, t.Progress*100) +
			"% " + fmt.Sprintf("状态: %s 错误: %+v", t.State, t.Err)
	}
}

func (t *TranslateTask) Type() TaskType {
	return TaskTypeTranslate
}

func (t *TranslateTask) GetState() TaskState {
	return t.State
}

func (g *TranslateTask) Failed(err error) {
	g.Err = err
	g.FailedTimes++
}

func (g *TranslateTask) GetFailedTimes() int {
	return g.FailedTimes
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
	name := strings.TrimRight(fileName, ext)
	t.DstFile = filepath.Join(absPath, fmt.Sprintf("%s_%s%s", name, t.To, ext))
	return
}

//字幕格式转换任务
type ConvertTask struct {
	SrcFile     string
	DstFile     string
	State       TaskState
	From        string
	To          string
	FailedTimes int
	Err         error
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
	if c.Err == nil {
		return fmt.Sprintf("[字幕格式转换] From: %s To: %s 源文件: %s 目标文件：%s  状态: %s",
			c.From, c.To, c.SrcFile, c.DstFile, c.State)
	} else {
		return fmt.Sprintf("[字幕格式转换] From: %s To: %s 源文件: %s 目标文件：%s  状态: %s 错误: %+v",
			c.From, c.To, c.SrcFile, c.DstFile, c.State, c.Err)
	}
}

func (c *ConvertTask) GetState() TaskState {
	return c.State
}

func (c *ConvertTask) Type() TaskType {
	return TaskTypeConvert
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
	name := strings.TrimRight(fileName, ext)
	c.DstFile = filepath.Join(absPath, fmt.Sprintf("%s.%s", name, c.To))
	return
}

func (g *ConvertTask) Failed(err error) {
	g.Err = err
	g.FailedTimes++
}

func (g *ConvertTask) GetFailedTimes() int {
	return g.FailedTimes
}

type GenerateStep int

const (
	GenerateStepUnknown = iota
	GenerateStepAudio
	GenerateStepOss
	GenerateStepRecognize
	GenerateStepGenerateSrt
)

func (g GenerateStep) String() string {
	switch g {
	case GenerateStepAudio:
		return "抽取音频"
	case GenerateStepOss:
		return "存储音频文件"
	case GenerateStepRecognize:
		return "语音识别"
	case GenerateStepGenerateSrt:
		return "生成字幕"
	default:
		return "未知步骤"
	}
}

type GenerateTask struct {
	SrcFile     string
	DstFile     string
	State       TaskState
	Step        GenerateStep
	ChannelId   int
	FailedTimes int
	Err         error
}

func NewGenerateTask(src string) *GenerateTask {
	return &GenerateTask{
		SrcFile: src,
		DstFile: "",
		State:   TaskStateInit,
		Step:    GenerateStepUnknown,
	}
}

func (g *GenerateTask) String() string {
	if g.Err == nil {
		return fmt.Sprintf("[字幕生成] 源文件: %s 目标文件: %s 步骤: %s 状态: %s",
			g.SrcFile, g.DstFile, g.Step, g.State)
	} else {
		return fmt.Sprintf("[字幕生成] 源文件: %s 目标文件: %s 步骤: %s 状态: %s 错误: %+v",
			g.SrcFile, g.DstFile, g.Step, g.State, g.Err)
	}
}

func (g *GenerateTask) Type() TaskType {
	return TaskTypeGenerate
}

func (g *GenerateTask) GetState() TaskState {
	return g.State
}

func (g *GenerateTask) Failed(err error) {
	g.Err = err
	g.FailedTimes++
}

func (g *GenerateTask) GetFailedTimes() int {
	return g.FailedTimes
}

type CleanTask struct {
	SrcFile     string
	DstFile     string
	State       TaskState
	Err         error
	Progress    float32
	FailedTimes int
}

func NewCleanTask(srcFile string) *CleanTask {
	return &CleanTask{
		SrcFile:     srcFile,
		DstFile:     "",
		State:       TaskStateInit,
		Err:         nil,
		Progress:    0,
		FailedTimes: 0,
	}
}

func (c *CleanTask) Init() (err error) {
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
	name := strings.TrimRight(fileName, ext)
	c.DstFile = filepath.Join(absPath, fmt.Sprintf("%s_%s%s", name, "clean", ext))

	return
}

func (t *CleanTask) String() string {
	if t.Err == nil {
		return fmt.Sprintf("源文件: %s 目标文件: %s 进度: %.2f", t.SrcFile, t.DstFile, t.Progress*100) +
			"% " + fmt.Sprintf("状态: %s", t.State)
	} else {
		return fmt.Sprintf("源文件: %s 目标文件: %s 进度: %.2f", t.SrcFile, t.DstFile, t.Progress*100) +
			"% " + fmt.Sprintf("状态: %s 错误: %+v", t.State, t.Err)
	}
}

func (t *CleanTask) Type() TaskType {
	return TaskTypeClean
}

func (t *CleanTask) GetState() TaskState {
	return t.State
}

func (g *CleanTask) Failed(err error) {
	g.Err = err
	g.FailedTimes++
}

func (g *CleanTask) GetFailedTimes() int {
	return g.FailedTimes
}
