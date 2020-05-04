package tencent

import (
	"encoding/json"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	api = "https://fanyi.qq.com/api/translate"
)

type Translator struct {
	qtk string
	qtv string
}

func NewTranslator(qtk, qtv string) *Translator {
	return &Translator{
		qtk: qtk,
		qtv: qtv,
	}
}
func (t *Translator) Do(src, from, to string) (dst string, err error) {
	var (
		params *url.Values
		ret    *response
	)
	params = &url.Values{}
	params.Add("source", from)
	params.Add("target", to)
	params.Add("sourceText", src)
	params.Add("sessionUuid", "translate_uuid"+strconv.FormatInt(time.Now().UnixNano()/1000/1000, 10))
	ret, err = t.call(params)
	if err != nil {
		return
	}

	if ret.Translate != nil && len(ret.Translate.Records) > 0 {
		dst = ret.Translate.Records[0].TargetText
	}
	return
}

func (t *Translator) call(params *url.Values) (ret *response, err error) {
	var (
		req     *http.Request
		rsp     *http.Response
		content []byte
	)
	url := api
	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	req.Header.Add("Cookie", fmt.Sprintf("qtv=%s; qtk=%s;", t.qtv, t.qtk))
	req.Header.Add("Origin", "http://fanyi.qq.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
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

	content, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	ret = new(response)
	json.Unmarshal(content, ret)
	return
}
