package srt

import (
	"bufio"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type state int

var (
	StateKnown    state = 0
	StateSequence state = 1
	StateTime     state = 2
	StateSubtitle state = 3
	StateEnd      state = 4
)

func Open(filePath string) (rets []*Srt, err error) {
	var (
		f          *os.File
		tmp        *Srt
		lineState  state
		line       string
		lineNumber int
	)

	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	lineNumber = 1
	lineState = StateSequence
	tmp = new(Srt)
	for {
		line, err = rd.ReadString('\n') //以'\n'为结束符读入一行
		if io.EOF == err {
			err = nil
			break
		}

		if err != nil {
			err = errors.WithMessage(status.ErrReadFileFailed, fmt.Sprintf("行号:%d", lineNumber))
			return
		}
		if line == "\n" {
			lineState = StateEnd
		}
		switch lineState {
		case StateEnd:
			rets = append(rets, tmp)
			tmp = new(Srt)
			lineState = StateSequence
		case StateSequence:
			tmp.Sequence, err = strconv.Atoi(strings.Trim(line, "\r\n"))
			if err != nil {
				err = errors.WithMessage(status.ErrSequence, fmt.Sprintf("文件: %s 行号: %d", filePath, lineNumber))
				return
			}
			lineState = StateTime
		case StateTime:
			lineState = StateSubtitle
			infos := strings.Split(line, " ")
			if len(infos) != 3 {
				err = errors.WithMessage(status.ErrTimeLine, fmt.Sprintf("文件：%s,行号：%d", filePath, lineNumber))
				return
			}
			tmp.Start = infos[0]
			tmp.End = infos[2]

		case StateSubtitle:
			tmp.Subtitle += line
		}

		lineNumber++
	}

	return
}

func WriteSrt(filePath string, src []*Srt) (err error) {
	var (
		absFilePath string
		dstFile     *os.File
	)
	sort.Sort(SrtSort(src))

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
		if itr.Subtitle == "" || itr.Subtitle == "\r\n" || itr.Subtitle == "\n" {
			continue
		}
		dstFile.WriteString(fmt.Sprintf("%d\r\n", itr.Sequence))
		dstFile.WriteString(fmt.Sprintf("%s --> %s\r\n", strings.Trim(itr.Start, "\r\n"),
			strings.Trim(itr.End, "\r\n")))
		dstFile.WriteString(strings.Trim(itr.Subtitle, "\r\n") + "\r\n")
		dstFile.WriteString("\r\n")
	}
	return
}
