package baidu

const (
	api = "http://api.fanyi.baidu.com/api/trans/vip/translate"
)

type transResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type response struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []transResult `json:"trans_result"`
	ErrorCode   string        `json:"error_code"`
	ErrorMsg    string        `json:"error_msg"`
}

const (
	OK = "52000"
)

var (
	ErrCode = map[string]string{
		"52000": "成功",
		"52001": "请求超时",
		"52002": "系统错误",
		"52003": "未授权用户",
		"54000": "必填参数为空",
		"54001": "签名错误",
		"54003": "访问频率受限",
		"54004": "账户余额不足",
		"54005": "长query请求频繁",
		"58000": "客户度IP非法",
		"58001": "译文语言方向不支持",
		"58002": "服务当前已关闭",
		"90107": "认证未通过或未生效",
	}
)
