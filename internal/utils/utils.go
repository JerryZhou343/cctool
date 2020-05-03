package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
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
