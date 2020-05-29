package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(level string, path, fileName string) {
	hook := &lumberjack.Logger{
		Filename:   path + fileName, // 日志文件路径
		MaxSize:    128,             // megabytes
		MaxBackups: 30,              // 最多保留300个备份
		MaxAge:     7,               // days
		Compress:   true,            // 是否压缩 disabled by default
	}
	l, err := logrus.ParseLevel(level)
	if err != nil {
		fmt.Printf("log level [%v]convert failed will be set to debug [%v]\n", level, err)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(l)
	}
	logrus.SetOutput(hook)
}
