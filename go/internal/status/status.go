package status

import "errors"

var (
	ErrPathError              = errors.New("路径不存在")
	ErrOpenFileFailed         = errors.New("打开文件失败")
	ErrReadFileFailed         = errors.New("读取文件错误")
	ErrSequence               = errors.New("序号转换错误")
	ErrTimeLine               = errors.New("时间行错误")
	ErrHttpCallFailed         = errors.New("http 请求失败")
	ErrTranslateFailed        = errors.New("翻译失败")
	ErrNotFoundConfig         = errors.New("没有找到配置文件")
	ErrInitTranslatorFailed   = errors.New("初始化翻译器失败")
	ErrCreatePathFailed       = errors.New("创建目标文件路径错误")
	ErrSubtitleNumberNoEnough = errors.New("字幕内容不够")
	ErrSourceFileNotEnough    = errors.New("源文件个数错误")
	ErrSourceFileMaxSize      = errors.New("源文件个数过多")
	ErrInitConfigFileFailed   = errors.New("加载配置文件错误")
	ErrDstFile                = errors.New("目标文件参数未填写")
	ErrFileNotExits           = errors.New("文件不存在")
	ErrFFmpegeCheckFailed     = errors.New("ffmpege 依赖检查失败,请安装")
	ErrConfigError            = errors.New("配置文件校验不通过")
	ErrSplitSentenceBug       = errors.New("Spilt Sentence bug")
)
