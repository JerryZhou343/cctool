package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/translate/common"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//https://api.fanyi.baidu.com/doc/21

type Translator struct {
	cfg *Config
}

func NewTranslator(cfg *Config) *Translator {
	ret := Translator{
		cfg: cfg,
	}
	return &ret
}

func (t *Translator) Do(src string, channel common.Channel) (dst string, err error) {
	var (
		params *url.Values
		ret    *response
	)
	params = &url.Values{}
	from, to := t.getFromTo(channel)
	salt := t.genSalt()
	params.Add("q", src)
	params.Add("appid", t.cfg.AppId)
	params.Add("salt", salt)
	params.Add("from", from)
	params.Add("to", to)
	params.Add("sign", t.genSign(src, salt))

	ret, err = t.call(params)
	if err != nil {
		return
	}

	if ret.ErrorCode != "0" && ret.ErrorCode != OK {
		err = errors.WithMessage(status.ErrHttpCallFailed, fmt.Sprintf("%+v", ErrCode[ret.ErrorCode]))
		return
	}
	if len(ret.TransResult) > 0{
		return ret.TransResult[0].Dst, nil
	}
	return "",status.ErrTranslateFailed
}

func (t *Translator) getFromTo(channel common.Channel) (string, string) {
	switch channel {
	case common.Channel_En2Zh:
		return "en", "zh"
	default:
		return "auto", "zh"
	}
}

func (t *Translator) genSign(src string, salt string) string {
	raw := t.cfg.AppId + src + salt + t.cfg.SecretKey
	return utils.Md5String(raw)
}

func (t *Translator) genSalt() string {
	return strconv.FormatInt(utils.GetIntRandomNumber(1000, 1000000), 10)
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

	content, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	log.Printf("%+v",string(content))
	ret = &response{}
	json.Unmarshal(content, ret)

	return
}
