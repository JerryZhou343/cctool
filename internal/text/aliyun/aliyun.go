package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ()

const (
	regexNumber = `^[\-0-9][0-9]*(.[0-9]+)?$`
)

type Speech struct {
	accessKeyId     string
	accessKeySecret string
	appKey          string
	wellKnownWord   map[string]string
	wellKnownNumber map[string]int
}

func NewSpeech(accessKeyId, accessKeySecret, appKey string,
	wellKnownNumber map[string]int, wellKnownWord map[string]string) text.ISpeech {
	return &Speech{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		appKey:          appKey,
		wellKnownNumber: wellKnownNumber,
		wellKnownWord:   wellKnownWord,
	}
}

func (s *Speech) Recognize(ctx context.Context, fileUri string) (sRet []*srt.Srt, wRet []*srt.Srt, err error) {
	var (
		taskId string
		rsp    *Response
		client *sdk.Client
	)

	client, err = sdk.NewClientWithAccessKey(REGION_ID, s.accessKeyId, s.accessKeySecret)
	if err != nil {
		return
	}
	logrus.Infof("aliyun speech add task %s",fileUri)
	taskId, err = s.sendTask(client, fileUri)
	if err != nil {
		return
	}

	rsp, err = s.queryResult(ctx, client, taskId)
	if err != nil || rsp == nil {
		return
	}
	if rsp.StatusText != STATUS_SUCCESS {
		err = errors.WithMessage(err, "recognize failed")
		return
	}
	logrus.Infof("aliyun speech recognize result [%d]",len(rsp.Result.Sentences))
	sRet, err = s.Sentence(0, rsp)
	wRet, err = s.BreakSentence(0, rsp)
	if err != nil {
		data, _ := json.Marshal(rsp)
		logrus.Errorf("break sentence failed [%+v]",err)
		_ = ioutil.WriteFile(fmt.Sprintf("log/dump_%d.json", time.Now().Unix()), data, os.ModePerm)
		err = nil
	}

	return
}

func (s *Speech) queryResult(ctx context.Context, client *sdk.Client, taskId string) (ret *Response, err error) {
	getRequest := requests.NewCommonRequest()
	getRequest.Domain = DOMAIN
	getRequest.Version = API_VERSION
	getRequest.Product = PRODUCT
	getRequest.ApiName = GET_REQUEST_ACTION
	getRequest.Method = "GET"
	getRequest.QueryParams[KEY_TASK_ID] = taskId
	var statusText = ""
	var getResponse *responses.CommonResponse
	for {
		select {
		case <-ctx.Done():
			return
		default:
			getResponse, err = client.ProcessCommonRequest(getRequest)
			if err != nil {
				return
			}
			if getResponse.GetHttpStatus() != 200 {
				err = errors.New(fmt.Sprintf("识别结果查询请求失败，Http错误码：%d", getResponse.GetHttpStatus()))
				break
			}
			ret = &Response{}
			json.Unmarshal(getResponse.GetHttpContentBytes(), ret)
			statusText = ret.StatusText
			if statusText == STATUS_RUNNING || statusText == STATUS_QUEUEING {
				time.Sleep(3 * time.Second)
				continue
			}
		}
		return
	}
}

func (s *Speech) sendTask(client *sdk.Client, URI string) (taskId string, err error) {
	postRequest := requests.NewCommonRequest()
	postRequest.Domain = DOMAIN
	postRequest.Version = API_VERSION
	postRequest.Product = PRODUCT
	postRequest.ApiName = POST_REQUEST_ACTION
	postRequest.Method = "POST"
	mapTask := make(map[string]string)
	mapTask[KEY_APP_KEY] = s.appKey
	mapTask[KEY_FILE_LINK] = URI
	// 新接入请使用4.0版本，已接入(默认2.0)如需维持现状，请注释掉该参数设置
	mapTask[KEY_VERSION] = "4.0"
	// 设置是否输出词信息，默认为false，开启时需要设置version为4.0
	mapTask[KEY_ENABLE_WORDS] = "true"
	mapTask[KEY_MAX_SINGLE_SEGMENT_TIME] = "1000"
	task, err := json.Marshal(mapTask)
	if err != nil {
		panic(err)
	}

	postRequest.FormParams[KEY_TASK] = string(task)
	postResponse, err := client.ProcessCommonRequest(postRequest)
	if err != nil {
		panic(err)
	}
	postResponseContent := postResponse.GetHttpContentString()
	if postResponse.GetHttpStatus() != 200 {
		err = errors.New(fmt.Sprintf("录音文件识别请求失败，Http错误码: %d", postResponse.GetHttpStatus()))
		return
	}
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		panic(err)
	}
	var statusText string = ""
	statusText = postMapResult[KEY_STATUS_TEXT].(string)
	if statusText == STATUS_SUCCESS {
		taskId = postMapResult[KEY_TASK_ID].(string)
	} else {
		err = errors.New(fmt.Sprintf("录音文件识别请求失败! %+s", statusText))
	}
	return
}

func (s *Speech) BreakSentence(channelId int, rsp *Response) (ret []*srt.Srt, err error) {
	var (
		newLine bool
		idx     int
	)
	//1. 重新断句
	idx = 0
	newLine = true
	tmpSrt := &srt.Srt{}
	for _, sentence := range rsp.Result.Sentences {
		//不是目标通道就过掉
		if sentence.ChannelId != channelId {
			continue
		}
		//1.1 按照空格切词
		words := strings.Split(sentence.Text, " ")
		//1.2 断句
		for _, word := range words {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}
			//句子结尾
			if strings.ContainsAny(word, ",.?!，。？！") {
				tmpSrt.Subtitle += " " + word
				newLine = true
				continue
			}
			//新句子开头
			if newLine == true {
				idx += 1
				tmpSrt = &srt.Srt{
					Sequence: idx,
					Subtitle: word,
				}
				ret = append(ret, tmpSrt)
				newLine = false
			} else { //句子中间
				tmpSrt.Subtitle += " " + word
			}
		}
	}

	re, _ := regexp.Compile(regexNumber)
	curIdx := 0
	for _, itr := range ret {
		sentenceWords := strings.Split(itr.Subtitle, " ")
		for swIdx, sw := range sentenceWords { //句子中的词
			sword := sw
			if strings.ContainsAny(sw, text.SentenceBreak) {
				sword = strings.TrimRight(sword, text.SentenceBreak)
			}

			sword = strings.ToLower(strings.TrimFunc(strings.TrimSpace(sword), func(r rune) bool {
				if strings.ContainsRune(text.SentenceBreak,r){
					return true
				}else{
					return false
				}
				return false
			}))

			numberFlag := false
			firstSetFlag := true
			for wIdx := curIdx; wIdx < len(rsp.Result.Words); wIdx++ {
				//更新curIdx
				if rsp.Result.Words[wIdx].ChannelId != channelId {
					continue
				}
				word := strings.ToLower(strings.TrimSpace(rsp.Result.Words[wIdx].Word))
				//fmt.Printf("%s:%s\n",word,sword)
				//当前单词是数字,并且句子中的词也是数字
				if v, ok := s.wellKnownNumber[word]; ok && re.Match([]byte(sword)) {
					numberFlag = true
					//如果首词是数字
					if swIdx == 0 && firstSetFlag {
						itr.Start = utils.MillisDurationConv(rsp.Result.Words[wIdx].BeginTime)
						firstSetFlag = false
					}
					itr.End = utils.MillisDurationConv(rsp.Result.Words[wIdx].EndTime)
					curIdx = wIdx + 1
					tmpNum, _ := strconv.Atoi(sword)
					//如果两个词相等，不再移动单词轴，句子轴和单词轴都移动
					if v == tmpNum {
						break
					}
					continue
				}

				//前面一个词是数字，但是现在这个词不是数字了
				if _, ok := s.wellKnownNumber[word]; !ok && numberFlag == true {
					numberFlag = false
					break
				}

				if !numberFlag && s.Equal(sword, word) {
					if swIdx == 0 {
						itr.Start = utils.MillisDurationConv(rsp.Result.Words[wIdx].BeginTime)
					}
					itr.End = utils.MillisDurationConv(rsp.Result.Words[wIdx].EndTime)
					curIdx = wIdx + 1
					break
				} else {
					err = errors.WithMessage(status.ErrSplitSentenceBug, fmt.Sprintf("sentence[%s] don't match dst:[%s] src[%s]", itr.Subtitle, sword, word))
					return
				}
			}
		}
	}
	return
}

//sw 句子中的词， w单词
func (s *Speech) Equal(sw, w string) bool {
	if v, ok := s.wellKnownWord[sw]; ok {
		if v == w {
			return true
		}
	}

	if sw == w {
		return true
	}

	return false
}

func (s *Speech) Sentence(channelId int, rsp *Response) (ret []*srt.Srt, err error) {
	var (
		idx = 0
	)
	for _, sentence := range rsp.Result.Sentences {
		if sentence.ChannelId != channelId {
			continue
		}
		idx += 1
		ret = append(ret, &srt.Srt{
			Sequence: idx,
			Start:    utils.MillisDurationConv(sentence.BeginTime),
			End:      utils.MillisDurationConv(sentence.EndTime),
			Subtitle: sentence.Text,
		})
	}

	return
}
