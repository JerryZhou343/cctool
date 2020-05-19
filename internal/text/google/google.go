package google

import (
	sp "cloud.google.com/go/speech/apiv1"
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"strings"
)

type Speech struct {
	credentialsFile string
	sampleRate      int32
}

func NewSpeech(credentialsFile string, sampleRate int32) text.ISpeech {
	return &Speech{
		credentialsFile: credentialsFile,
		sampleRate:      sampleRate,
	}
}

func (s *Speech) Recognize(ctx context.Context, fileURI string) (ret []*srt.Srt, err error) {
	var (
		client *sp.Client
	)
	//准备数据
	client, err = sp.NewClient(context.Background(), option.WithCredentialsFile(s.credentialsFile))
	if err != nil {
		err = errors.Wrap(err, "创建google客户端失败")
		return
	}

	//
	audio := &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Uri{
			Uri: fileURI,
		},
	}

	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                            speechpb.RecognitionConfig_ENCODING_UNSPECIFIED,
			SampleRateHertz:                     s.sampleRate,
			EnableSeparateRecognitionPerChannel: true,
			LanguageCode:                        "en-US",
			SpeechContexts:                      nil, //todo:识别场景
			EnableWordTimeOffsets:               true,
			EnableAutomaticPunctuation:          true,
			AudioChannelCount:                   8,
		},
		Audio: audio,
	}
	var (
		rsp       *speechpb.LongRunningRecognizeResponse
		operation *speech.LongRunningRecognizeOperation
	)

	operation, err = client.LongRunningRecognize(ctx, req)
	if err != nil {
		return
	}
	rsp, err = operation.Wait(ctx)
	if err != nil {
		return
	}
	var idx = 0
	var newLine bool
	var tmpSrt *srt.Srt
	newLine = true
	for _, result := range rsp.Results {
		if result.ChannelTag != 1 {
			continue
		}
		for _, itr := range result.Alternatives {
			for _, word := range itr.Words {
				//句子结尾
				if strings.ContainsAny(word.Word, ",.?!，。？！") {
					tmpSrt.End = utils.DurationConv(word.EndTime)
					tmpSrt.Subtitle += " " + word.Word
					newLine = true
					continue
				}
				//新句子开头
				if newLine == true {
					idx += 1
					tmpSrt = &srt.Srt{
						Sequence: idx,
						Start:    utils.DurationConv(word.StartTime),
						End:      utils.DurationConv(word.EndTime),
						Subtitle: word.Word,
					}
					ret = append(ret, tmpSrt)
					newLine = false
				} else { //句子中间
					tmpSrt.End = utils.DurationConv(word.EndTime)
					tmpSrt.Subtitle += " " + word.Word
				}
			}
		}
	}

	return
}
