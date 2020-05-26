package tencent

type transResult struct {
	SourceText string `json:"sourceText"`
	TargetText string `json:"targetText"`
	TraceId    string `json:"traceId"`
}

type response struct {
	SessionUuid string     `json:"sessionUuid"`
	Translate   *translate `json:"translate"`
}
type translate struct {
	ErrCode int            `json:"errCode"`
	ErrMsg  string         `json:"errMsg"`
	Source  string         `json:"source"`
	Target  string         `json:"target"`
	Records []*transResult `json:"records"`
}
