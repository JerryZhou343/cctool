package convert

import (
	"github.com/JerryZhou343/cctool/internal/bcc"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/utils"
	"sort"
)

func BCC2SRT(src *bcc.BCC) (ret []*srt.Srt) {
	sort.Sort(bcc.SortSubtitle(src.Body))
	for idx, itr := range src.Body {
		tmp := srt.Srt{}
		tmp.Sequence = idx + 1
		tmp.Start = utils.MillisDurationConv(int64(itr.From * 1000))
		tmp.End = utils.MillisDurationConv(int64(itr.To * 1000))
		tmp.Subtitle = itr.Content
		ret = append(ret, &tmp)
	}
	return
}
