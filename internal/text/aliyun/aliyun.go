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

func (s *Speech) Recognize(ctx context.Context, fileUri string, channelId int) (ret []*srt.Srt, err error) {
	client, err := sdk.NewClientWithAccessKey(REGION_ID, s.accessKeyId, s.accessKeySecret)
	if err != nil {
		panic(err)
	}

	postRequest := requests.NewCommonRequest()
	postRequest.Domain = DOMAIN
	postRequest.Version = API_VERSION
	postRequest.Product = PRODUCT
	postRequest.ApiName = POST_REQUEST_ACTION
	postRequest.Method = "POST"
	mapTask := make(map[string]string)
	mapTask[KEY_APP_KEY] = s.appKey
	mapTask[KEY_FILE_LINK] = fileUri
	// 新接入请使用4.0版本，已接入(默认2.0)如需维持现状，请注释掉该参数设置
	mapTask[KEY_VERSION] = "4.0"
	// 设置是否输出词信息，默认为false，开启时需要设置version为4.0
	mapTask[KEY_ENABLE_WORDS] = "false"
	mapTask[KEY_MAX_SINGLE_SEGMENT_TIME] = "900"
	//mapTask[KEY_ENABLE_DISFLUENCY] = "true"
	//mapTask[KEY_ENABLE_UNIFY_POST] = "true"
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
	fmt.Println(postResponseContent)
	if postResponse.GetHttpStatus() != 200 {
		fmt.Println("录音文件识别请求失败，Http错误码: ", postResponse.GetHttpStatus())
		return
	}
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		panic(err)
	}
	var taskId string = ""
	var statusText string = ""
	statusText = postMapResult[KEY_STATUS_TEXT].(string)
	if statusText == STATUS_SUCCESS {
		fmt.Println("录音文件识别请求成功响应!")
		taskId = postMapResult[KEY_TASK_ID].(string)
	} else {
		fmt.Println("录音文件识别请求失败!")
		return
	}

	getRequest := requests.NewCommonRequest()
	getRequest.Domain = DOMAIN
	getRequest.Version = API_VERSION
	getRequest.Product = PRODUCT
	getRequest.ApiName = GET_REQUEST_ACTION
	getRequest.Method = "GET"
	getRequest.QueryParams[KEY_TASK_ID] = taskId
	statusText = ""

	for {
		select {
		case <-ctx.Done():
			return
		default:
			getResponse, err := client.ProcessCommonRequest(getRequest)
			if err != nil {
				return ret, err
			}
			//getResponseContent := getResponse.GetHttpContentString()
			//fmt.Println("识别查询结果：", getResponseContent)
			if getResponse.GetHttpStatus() != 200 {
				fmt.Println("识别结果查询请求失败，Http错误码：", getResponse.GetHttpStatus())
				break
			}
			var rsp Response
			json.Unmarshal(getResponse.GetHttpContentBytes(), &rsp)
			statusText = rsp.StatusText
			if statusText == STATUS_RUNNING || statusText == STATUS_QUEUEING {
				time.Sleep(3 * time.Second)
				continue
			} else {
				if statusText == STATUS_SUCCESS {
					idx := 0
					for _, itr := range rsp.Result.Sentences {
						if itr.ChannelId == channelId {
							idx += 1
							tmpSrt := &srt.Srt{
								Sequence: idx,
								Start:    utils.MillisDurationConv(itr.BeginTime),
								End:      utils.MillisDurationConv(itr.EndTime),
								Subtitle: itr.Text,
							}
							ret = append(ret, tmpSrt)
						}
					}
				}
				break
			}
		}

	}
	return
}
