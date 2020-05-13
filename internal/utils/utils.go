package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/ptypes/duration"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetIntRandomNumber(min int64, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func CheckFileExist(path string) (flag bool) {
	var (
		err     error
		absFile string
	)
	absFile, err = filepath.Abs(path)
	if err != nil {
		return
	}

	_, err = os.Stat(absFile)
	if err == nil {
		flag = true
		return
	}

	if os.IsNotExist(err) {
		flag = false
	}
	return
}

func MillisDurationConv(d int64) (ret string) {
	var (
		millis  int64
		seconds int64
		mins    int64
		hours   int64
	)
	//微秒
	millis = (d % 1000)
	//秒
	seconds = (d / 1000)
	if seconds > 59 {
		mins = (d / 1000) / 60
		seconds = seconds % 60
	}
	//分
	if mins > 59 {
		hours = (d / 1000) / 3600
		mins = mins % 60
	}

	strHours := fmt.Sprintf("%02d", hours)
	strMins := fmt.Sprintf("%02d", mins)
	strSeconds := fmt.Sprintf("%02d", seconds)
	strMills := fmt.Sprintf("%03d", millis)

	return strHours + ":" + strMins + ":" + strSeconds + "," + strMills
}

func DurationConv(src *duration.Duration) string {
	var (
		millis  int64
		seconds int64
		mins    int64
		hours   int64
	)
	millis = int64(src.Nanos) / 1000 / 1000

	d := src.Seconds
	seconds = d
	if seconds > 59 {
		mins = (d / 1000) / 60
		seconds = seconds % 60
	}
	//分
	if mins > 59 {
		hours = (d / 1000) / 3600
		mins = mins % 60
	}

	strHours := fmt.Sprintf("%02d", hours)
	strMins := fmt.Sprintf("%02d", mins)
	strSeconds := fmt.Sprintf("%02d", seconds)
	strMills := fmt.Sprintf("%03d", millis)

	return strHours + ":" + strMins + ":" + strSeconds + "," + strMills
}
