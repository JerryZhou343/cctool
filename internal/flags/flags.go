package flags

var (
	SrcFiles       []string
	From           string
	To             string
	TransTool      string
	Merge          bool
	MergeStrategy  string
	DstFile        string
	AudioChannelId int
)

//translate tool
var (
	TransTool_Baidu   = "baidu"
	TransTool_Tencent = "tencent"
	TransTool_Google  = "google"
)

//merge strategy
var (
	StrategySequence = "seq"
	StrategyTimeline = "timeline"
)
