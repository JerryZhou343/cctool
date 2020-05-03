package status

import "errors"

var (
	ErrPathError      = errors.New("路径不存在")
	ErrOpenFileFailed = errors.New("打开文件失败")
	ErrReadFileFailed = errors.New("读取文件错误")
	ErrSequence       = errors.New("序号转换错误")
	ErrTimeLine       = errors.New("时间行错误")
	ErrHttpCallFailed = errors.New("http 请求失败")
	ErrTranslateFailed = errors.New("翻译失败")
)
