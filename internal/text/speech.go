package text

import "github.com/JerryZhou343/cctool/internal/srt"

type ISpeech interface {
	Recognize(fileURI string, channelId int) ([]*srt.Srt, error)
}
