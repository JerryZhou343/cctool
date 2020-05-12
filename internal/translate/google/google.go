package google

import (
	"encoding/json"
	"github.com/JerryZhou343/cctool/internal/status"
	gt "github.com/kyai/google-translate-tk"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	api = "https://translate.google.cn/translate_a/single"
)

type Translator struct {
	ttk string
}

func NewTranslator() *Translator {
	tkk, _ := gt.GetTKK()
	ret := Translator{
		ttk: tkk,
	}
	return &ret
}

func (t *Translator) Do(src, from, to string) (dst string, err error) {
	var (
		params *url.Values
		ret    *response
	)
	params = &url.Values{}
	params.Add("q", src)
	params.Add("client", "t")
	params.Add("sl", from)
	params.Add("tl", to)
	params.Add("hl", "ah-CN")
	params.Add("dt", "at")
	params.Add("dt", "bd")
	params.Add("dt", "ex")
	params.Add("dt", "ld")
	params.Add("dt", "md")
	params.Add("dt", "qca")
	params.Add("dt", "rw")
	params.Add("dt", "rm")
	params.Add("dt", "ss")
	params.Add("dt", "t")
	params.Add("ie", "UTF-8")
	params.Add("oe", "UTF-8")
	params.Add("source", "btn")
	params.Add("ssel", "0")
	params.Add("tsel", "0")
	params.Add("kc", "0")
	params.Add("tk", gt.GetTK(src, t.ttk))

	ret, err = t.call(params)
	if err != nil {
		return
	}
	dst = ret.Dst
	return
}

func (t *Translator) call(params *url.Values) (ret *response, err error) {
	var (
		req     *http.Request
		rsp     *http.Response
		content []byte
	)
	url := api + "?" + params.Encode()
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.WithMessagef(status.ErrHttpCallFailed, "%+v", err)
		return
	}

	http.DefaultClient.Timeout = time.Second * 5
	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		err = errors.WithMessagef(status.ErrHttpCallFailed, "%+v", err)
		return
	}
	defer rsp.Body.Close()

	content, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	//fmt.Println(string(content))
	var arr []interface{}
	json.Unmarshal(content, &arr)
	ret = new(response)
	for _, itr := range arr[0].([]interface{}) {
		if v, ok := itr.([]interface{})[0].(string); ok {
			ret.Dst += v
		}
	}
	return
}
