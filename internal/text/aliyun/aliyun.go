package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/srt"
	"github.com/JerryZhou343/cctool/internal/text"
	"github.com/JerryZhou343/cctool/internal/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type Speech struct {
	accessKeyId     string
	accessKeySecret string
	appKey          string
}

func NewSpeech(accessKeyId, accessKeySecret, appKey string) text.ISpeech {
	return &Speech{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		appKey:          appKey,
	}
}

func (s *Speech) Recognize(ctx context.Context, fileUri string) (ret []*srt.Srt, err error) {
	var (
		taskId string
		rsp    *Response
	)

	client, err := sdk.NewClientWithAccessKey(REGION_ID, s.accessKeyId, s.accessKeySecret)
	if err != nil {
		return
	}

	taskId, err = s.sendTask(client, fileUri)
	if err != nil {
		return
	}

	rsp, err = s.queryResult(ctx, client, taskId)
	if err != nil || rsp == nil {
		return
	}
	if rsp.StatusText != STATUS_SUCCESS {
		err = errors.New("recognize failed")
		return
	}
	ret = s.BreakSentence(0, rsp)

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
				fmt.Println("识别结果查询请求失败，Http错误码：", getResponse.GetHttpStatus())
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

func (s *Speech) BreakSentence(channelId int, rsp *Response) (ret []*srt.Srt) {
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

	curIdx := 0
	for _, itr := range ret {
		sentenceWords := strings.Split(itr.Subtitle, " ")
		for swIdx, sw := range sentenceWords {
			sword := sw
			if strings.ContainsAny(sw, text.SentenceBreak) {
				sword = strings.TrimRight(sword, text.SentenceBreak)
			}
			for wIdx := curIdx; wIdx < len(rsp.Result.Words); wIdx++ {
				//更新curIdx
				curIdx = wIdx + 1
				if rsp.Result.Words[wIdx].ChannelId != channelId {
					continue
				}
				if s.equal(sword, rsp.Result.Words[wIdx].Word) {
					//fmt.Println(rsp.Result.Words[wIdx].Word)
					if swIdx == 0 {
						itr.Start = utils.MillisDurationConv(rsp.Result.Words[wIdx].BeginTime)
					}
					itr.End = utils.MillisDurationConv(rsp.Result.Words[wIdx].EndTime)
					break
				} else {
					//fmt.Println(sword, ":", rsp.Result.Words[wIdx].Word, ": wIdx", wIdx, "sIdx", sIdx)
					//fmt.Println("dont't match")
					continue
					//return
				}
			}
		}
	}
	return
}

var (
	number = map[string]string{
		"1":    "one",
		"2":    "two",
		"3":    "three",
		"4":    "four",
		"5":    "five",
		"6":    "six",
		"7":    "seven",
		"8":    "eight",
		"9":    "nine",
		"0":    "zero",
		"x00":  "hundred",
		"x000": "thousand",
	}
)

//src 为源句子， dst 源为词
func (s *Speech) equal(src, dst string) bool {
	src = strings.ToLower(strings.TrimSpace(src))
	dst = strings.ToLower(strings.TrimSpace(dst))
	if src == dst {
		return true
	}

	//todo: 更换策略
	//if v, ok := number[src]; ok {
	//	if v == dst {
	//		return true
	//	}
	//}
	//
	_, err := strconv.ParseFloat(src, 10)
	if err == nil {
		return true
	}

	for _, v := range number {
		if v == dst {
			return true
		}
	}

	return false
}
