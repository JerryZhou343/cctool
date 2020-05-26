package tencent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetTVTK() (tv, tk string, err error) {
	resp, err := http.Get(api)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	data := string(body)
	fmt.Println(data)

	rtk := regexp.MustCompile(`qtk\s=\s"[\S]*"`)
	rtv := regexp.MustCompile(`qtv\s=\s"[\S]*"`)

	if rtk.MatchString(data) {
		tk = rtk.FindStringSubmatch(data)[1]
	}

	if rtv.MatchString(data) {
		tk = rtv.FindStringSubmatch(data)[1]
	}
	return
}
