package text

import (
	"context"
	"github.com/JerryZhou343/cctool/go/internal/srt"
)

type ISpeech interface {
	Recognize(ctx context.Context, fileURI string) ([]*srt.Srt, error)
}
