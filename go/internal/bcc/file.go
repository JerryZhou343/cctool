package bcc

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func Open(file string) (ret *BCC, err error) {
	var (
		absPath string
		content []byte
	)
	absPath, err = filepath.Abs(file)
	if err != nil {
		return
	}

	content, err = ioutil.ReadFile(absPath)
	if err != nil {
		return
	}

	ret = &BCC{}

	err = json.Unmarshal(content, ret)

	return
}
