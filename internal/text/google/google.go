package google

import (
	sp "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io/ioutil"
	"path/filepath"
)

type Speech struct {
	credentialsFile string
}

func NewSpeech(credentialsFile string) text.ISpeech {
	return &Speech{
		credentialsFile: credentialsFile,
	}
}

func (s *Speech) Recognize(srcFile string) (ret []*srt.Srt, err error) {
	var (
		absFile string
		content []byte
		client  *sp.Client
	)
	//准备数据
	absFile, err = filepath.Abs(srcFile)
	if err != nil {
		return
	}
	content, err = ioutil.ReadFile(absFile)
	if err != nil {
		return
	}

	client, err = sp.NewClient(context.Background(), option.WithCredentialsFile(s.credentialsFile))
	if err != nil {
		err = errors.Wrap(err, "创建google客户端失败")
		return
	}

	//
	audio := &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Content{
			Content: content,
		},
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                            speechpb.RecognitionConfig_SPEEX_WITH_HEADER_BYTE,
			SampleRateHertz:                     int32(conf.G_Config.SampleRate),
			EnableSeparateRecognitionPerChannel: true,
			LanguageCode:                        "en",
			SpeechContexts:                      nil, //todo:识别场景
			EnableWordTimeOffsets:               false,
			EnableAutomaticPunctuation:          true,
		},
		Audio: audio,
	}
	var (
		rsp *speechpb.RecognizeResponse
	)

	rsp, err = client.Recognize(context.Background(), req)
	if err == nil {
		fmt.Printf("result %d", len(rsp.Results))
	}
	return
}
