package aliyun

// 地域ID，常量内容，请勿改变
const REGION_ID string = "cn-shanghai"
const ENDPOINT_NAME string = "cn-shanghai"
const PRODUCT string = "nls-filetrans"
const DOMAIN string = "filetrans.cn-shanghai.aliyuncs.com"
const API_VERSION string = "2018-08-17"
const POST_REQUEST_ACTION string = "SubmitTask"
const GET_REQUEST_ACTION string = "GetTaskResult"

// 请求参数key
const KEY_APP_KEY string = "appkey"
const KEY_FILE_LINK string = "file_link"
const KEY_VERSION string = "version"
const KEY_ENABLE_WORDS string = "enable_words"
const KEY_MAX_SINGLE_SEGMENT_TIME  = "max_single_segment_time"
const KEY_ENABLE_UNIFY_POST = "enable_unify_post"
const KEY_ENABLE_DISFLUENCY = "enable_disfluency"

// 响应参数key
const KEY_TASK string = "Task"
const KEY_TASK_ID string = "TaskId"
const KEY_STATUS_TEXT string = "StatusText"
const KEY_RESULT string = "Result"

// 状态值
const STATUS_SUCCESS string = "SUCCESS"
const STATUS_RUNNING string = "RUNNING"
const STATUS_QUEUEING string = "QUEUEING"

type Result struct {
	EndTime         int64  `json:"EndTime"`
	SilenceDuration int    `json:"SilenceDuration"`
	BeginTime       int64  `json:"BeginTime"`
	Text            string `json:"Text"`
	ChannelId       int    `json:"ChannelId"`
	SpeechRate      int    `json:"SpeechRate"`
	EmotionValue    int    `json:"EmotionValue"`
}

type SentencesResult struct {
	Sentences []*Result `json:"Sentences"`
}

type Response struct {
	TaskId      string           `json:"TaskId"`
	RequestId   string           `json:"RequestId"`
	StatusText  string           `json:"StatusText"`
	BizDuration int64            `json:"BizDuration"`
	SolveTime   int64            `json:"SolveTime"`
	StatusCode  int32            `json:"StatusCode"`
	Result      *SentencesResult `json:"Result"`
}
