package conf

import (
	"github.com/JerryZhou343/cctool/internal/status"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	G_Config Config
)

type ApiConf struct {
	AppId     string `yaml:"app_id"`
	SecretKey string `yaml:"secret_key"`
	Interval  int64  `yaml:"interval"`
}

type TencentConf struct {
	Interval int    `yaml:"interval"`
	Qtv      string `yaml:"qtv"`
	Qtk      string `yaml:"qtk"`
}

type Config struct {
	Baidu   ApiConf     `yaml:"baidu"`
	Google  ApiConf     `yaml:"google"`
	Tencent TencentConf `yaml:"tencent"`
}

func Init() error {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return status.ErrNotFoundConfig
	}
	err = yaml.Unmarshal(f, &G_Config)
	log.Printf("%s: %+v \n", string(f), err)
	log.Printf("%+v", G_Config)
	return err
}
