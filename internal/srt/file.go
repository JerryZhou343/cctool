package srt

import (
	"bufio"
	"fmt"
	"github.com/JerryZhou343/ClosedCaption/internal/status"
	"github.com/pkg/errors"
	"io"
	"os"
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
			tmp.Sequence, err = strconv.Atoi(strings.Trim(line, "\n"))
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
			tmp.start = infos[0]
			tmp.end = infos[2]

		case StateSubtitle:
			tmp.subtitle += line
		}

		lineNumber++
	}

	return
}
