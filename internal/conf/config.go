package conf

import (
	"github.com/JerryZhou343/cctool/internal/status"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	G_Config Config
)

type ApiConf struct {
	AppId     string `yaml:"app_id"`
	SecretKey string `yaml:"secret_key"`
	Interval  int64  `yaml:"interval"`
}

func (a *ApiConf) Check() bool {
	if a.AppId == "" || a.SecretKey == "" {
		return false
	}
	if a.Interval == 0 {
		a.Interval = 1000
	}
	return true
}

type TencentConf struct {
	Interval int    `yaml:"interval"`
	Qtv      string `yaml:"qtv"`
	Qtk      string `yaml:"qtk"`
}

func (t *TencentConf) Check() bool {
	if t.Qtk == "" || t.Qtv == "" {
		return false
	}
	if t.Interval == 0 {
		t.Interval = 1000
	}
	return true
}

type AliYunConf struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"accessKey_secret"`
	AppKey          string `yaml:"app_key"`
	OssEndpoint     string `yaml:"oss_endpoint"`
	BucketName      string `yaml:"bucket_name"`
	BucketDomain    string `yaml:"bucket_domain"`
}

func (a *AliYunConf) Check() bool {
	if a.AccessKeyId == "" ||
		a.AccessKeySecret == "" ||
		a.AppKey == "" ||
		a.OssEndpoint == "" ||
		a.BucketDomain == "" ||
		a.BucketName == "" {
		return false
	}
	return true
}

type Config struct {
	Baidu          ApiConf     `yaml:"baidu"`
	Google         ApiConf     `yaml:"google"`
	Tencent        TencentConf `yaml:"tencent"`
	Aliyun         AliYunConf  `yaml:"aliyun"`
	SampleRate     int
	AudioCachePath string   `yaml:"audio_cache_path"`
	SrtPath        string   `yaml:"srt_path"`
	TransTools     []string `yaml:"translate_tool"`
}

func Load() (err error) {
	var (
		f        []byte
		currPath string
	)
	f, err = ioutil.ReadFile("config.yaml")
	if err != nil {
		return status.ErrNotFoundConfig
	}
	err = yaml.Unmarshal(f, &G_Config)
	G_Config.SampleRate = 16000
	currPath, err = os.Getwd()
	if G_Config.AudioCachePath == "" {
		G_Config.AudioCachePath = filepath.Join(currPath, "audio")
	}

	if G_Config.SrtPath == "" {
		G_Config.SrtPath = filepath.Join(currPath, "srt")
	}
	G_Config.AudioCachePath, err = filepath.Abs(G_Config.AudioCachePath)
	G_Config.SrtPath, err = filepath.Abs(G_Config.SrtPath)
	return err
}
