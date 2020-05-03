package status

import "errors"

var (
	ErrPathError            = errors.New("路径不存在")
	ErrOpenFileFailed       = errors.New("打开文件失败")
	ErrReadFileFailed       = errors.New("读取文件错误")
	ErrSequence             = errors.New("序号转换错误")
	ErrTimeLine             = errors.New("时间行错误")
	ErrHttpCallFailed       = errors.New("http 请求失败")
	ErrTranslateFailed      = errors.New("翻译失败")
	ErrNotFoundConfig       = errors.New("没有找到配置文件")
	ErrInitTranslatorFailed = errors.New("初始化翻译器失败")
	ErrCreatePathFailed     = errors.New("创建目标文件路径错误")
)
