package merge

import (
	"fmt"
	"github.com/JerryZhou343/cctool/go/internal/flags"
	"github.com/JerryZhou343/cctool/go/internal/srt"
	"github.com/JerryZhou343/cctool/go/internal/status"
	"github.com/pkg/errors"
	"path/filepath"
	"sort"
)

type Merge struct {
	srcSrt []map[int]*srt.Srt
}

func NewMerge() *Merge {
	return &Merge{
		srcSrt: []map[int]*srt.Srt{},
	}
}

func (e *Merge) open(srcFilePaths ...string) (err error) {
	var (
		tmpAbsSrcFile string
		rets          []*srt.Srt
	)
	for _, itr := range srcFilePaths {
		tmpAbsSrcFile, err = filepath.Abs(itr)
		if err != nil {
			err = errors.WithMessage(status.ErrPathError, fmt.Sprintf("%s", itr))
		}
		rets, err = srt.Open(tmpAbsSrcFile)
		if err != nil {
			return
		}
		tmp := map[int]*srt.Srt{}
		for _, itr := range rets {
			tmp[itr.Sequence] = itr
		}
		e.srcSrt = append(e.srcSrt, tmp)
	}

	if len(e.srcSrt) <= 1 {
		return status.ErrSubtitleNumberNoEnough
	}
	return
}

func (e *Merge) Merge(strategy string, dstFilePath string, srcFilePaths ...string) (err error) {
	//参数校验
	if len(srcFilePaths) <= 1 {
		return status.ErrSourceFileNotEnough
	} else if len(srcFilePaths) > 2 {
		return status.ErrSourceFileMaxSize
	}

	//打开文件
	err = e.open(srcFilePaths...)
	if err != nil {
		return
	}

	switch strategy {
	case flags.StrategySequence:
		err = e.single(dstFilePath)
	case flags.StrategyTimeline:
		err = e.mix(dstFilePath)
	}
	return
}

//按照序列号合并
func (e *Merge) single(dstFilePath string) (err error) {
	base := e.srcSrt[0]
	for i := 1; i < len(e.srcSrt); i++ {
		for _, itr := range e.srcSrt[i] {
			if v, ok := base[itr.Sequence]; ok {
				v.Subtitle += fmt.Sprintf("%s\r\n", itr.Subtitle)
			}
		}
	}
	src := []*srt.Srt{}
	for _, itr := range base {
		src = append(src, itr)
	}
	err = srt.WriteSrt(dstFilePath, src)
	return
}

//时间轴起点相同的区间进行合并，重新生成序号
func (e *Merge) mix(dstFilePath string) (err error) {
	//数据准备
	src1Set := e.srcSrt[0]
	src2Set := e.srcSrt[1]

	src1Slice := []*srt.Srt{}
	src2Slice := []*srt.Srt{}
	src2StartTimeSet := map[string]*srt.Srt{}
	src2EndTimeSet := map[string]*srt.Srt{}

	for _, itr := range src1Set {
		src1Slice = append(src1Slice, itr)
	}
	sort.Sort(srt.SrtSort(src1Slice))

	for _, itr := range src2Set {
		src2Slice = append(src2Slice, itr)
		src2EndTimeSet[itr.End] = itr
		src2StartTimeSet[itr.Start] = itr
	}
	sort.Sort(srt.SrtSort(src2Slice))

	for _, itr := range src1Slice {
		//查找时间轴闭区间集合
		startSeq := 0
		endSeq := 0
		if v, ok := src2StartTimeSet[itr.Start]; ok {
			startSeq = v.Sequence
		}

		if v, ok := src2EndTimeSet[itr.End]; ok {
			endSeq = v.Sequence
		}

		for _, item := range src2Slice {
			if item.Sequence >= startSeq && item.Sequence <= endSeq {
				itr.Subtitle += item.Subtitle
			}
		}
	}
	srt.WriteSrt(dstFilePath, src1Slice)
	return
}
