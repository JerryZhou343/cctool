package app

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/translate"
	"github.com/JerryZhou343/cctool/internal/translate/baidu"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Application struct {
	translator translate.Translate
	interval   time.Duration
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Translate() (err error) {
	var (
		ret []*srt.Srt
	)
	switch flags.TransTool {
	case flags.TransTool_Baidu:
		a.translator = baidu.NewTranslator(&conf.G_Config.Baidu)
		a.interval = time.Millisecond * time.Duration(conf.G_Config.Baidu.Interval)
	default:
		return status.ErrInitTranslatorFailed
	}

	for _, itr := range flags.SrcFiles {
		ret, err = srt.Open(itr)
		if err != nil {
			return err
		}
		err = a.translate(itr, ret)
		if err != nil {
			return
		}
	}

	return
}

func (a *Application) translate(srcPath string, src []*srt.Srt) (err error) {
	var (
		tmpResult   string
		absFilePath string
		ret         map[int]string
	)
	absFilePath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	absPath := filepath.Dir(absFilePath)
	fileName := filepath.Base(absFilePath)
	dstFile := filepath.Join(absPath, flags.To+"_"+fileName)


	ret = make(map[int]string)
	for _, itr := range src {
		tmpResult, err = a.translator.Do(itr.Subtitle, flags.From, flags.To)
		if err != nil {
			return
		}
		ret[itr.Sequence] = tmpResult
		time.Sleep(a.interval)
	}

	if flags.Merge {
		for _, itr := range src {
			if v, ok := ret[itr.Sequence]; ok {
				itr.Subtitle = v + "\n" + itr.Subtitle
			}
		}
		a.WriteSrt(dstFile, src)
	} else {
		for _, itr := range src {
			if v, ok := ret[itr.Sequence]; ok {
				itr.Subtitle = v
			}
		}
		a.WriteSrt(dstFile, src)
	}
	return
}

func (a *Application) WriteSrt(filePath string, src []*srt.Srt) (err error) {
	var (
		absFilePath string
		dstFile     *os.File
	)
	sort.Sort(srt.SrtSort(src))

	absFilePath, err = filepath.Abs(filePath)
	if err != nil {
		err = errors.WithMessage(status.ErrPathError, fmt.Sprintf("%s", filePath))
		return
	}

	absPath := filepath.Dir(absFilePath)
	_, err = os.Stat(absPath)
	if err != nil {
		err = os.Mkdir(absPath, os.ModePerm)
		if err != nil {
			return errors.WithMessage(status.ErrCreatePathFailed, fmt.Sprintf("路径: %s", absPath))
		}
	}

	dstFile, err = os.OpenFile(absFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.WithMessage(status.ErrOpenFileFailed, fmt.Sprintf("文件：%s", filePath))
	}
	defer dstFile.Close()

	for _, itr := range src {
		dstFile.WriteString(fmt.Sprintf("%d\r\n", itr.Sequence))
		dstFile.WriteString(fmt.Sprintf("%s --> %s", itr.Start, itr.End))
		dstFile.WriteString(itr.Subtitle)
		dstFile.WriteString("\r\n")
	}
	return
}
