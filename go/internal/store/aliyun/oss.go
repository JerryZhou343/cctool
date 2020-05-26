package aliyun

import (
	"github.com/JerryZhou343/cctool/go/internal/status"
	"github.com/JerryZhou343/cctool/go/internal/store"
	"github.com/JerryZhou343/cctool/go/internal/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AliyunOSS struct {
	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucketName      string
	bucketDomain    string
}

func NewAliyunOSS(endPoint, accessKeyId, accessKeySecret, bucketName, bucketDomain string) store.Store {
	return &AliyunOSS{
		endpoint:        endPoint,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		bucketName:      bucketName,
		bucketDomain:    bucketDomain,
	}
}

func (a *AliyunOSS) GetListBuckets() (ret []string, err error) {
	var (
		client *oss.Client
		rsp    oss.ListBucketsResult
	)

	client, err = oss.New(a.endpoint, a.accessKeyId, a.accessKeySecret)
	if err != nil {
		return nil, err
	}

	rsp, err = client.ListBuckets()
	if err != nil {
		return nil, err
	}

	for _, bucket := range rsp.Buckets {
		ret = append(ret, bucket.Name)
	}
	return
}

func (a *AliyunOSS) UploadFile(srcFilePath string) (uri string, objName string, err error) {
	var (
		client *oss.Client
		bucket *oss.Bucket
	)
	if !utils.CheckFileExist(srcFilePath) {
		err = status.ErrFileNotExits
		return
	}
	fileName := filepath.Base(srcFilePath)
	client, err = oss.New(a.endpoint, a.accessKeyId, a.accessKeySecret)
	if err != nil {
		return
	}
	bucket, err = client.Bucket(a.bucketName)
	if err != nil {
		return
	}

	//分日期存储
	date := time.Now()
	year := date.Year()
	month := date.Month()
	day := date.Day()
	objName = strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/" + fileName

	err = bucket.PutObjectFromFile(objName, srcFilePath)
	if err != nil {
		return
	}

	return a.GetObjectFileUrl(objName), objName, nil
}

func (a *AliyunOSS) DeleteFile(uri string) error {
	client, err := oss.New(a.endpoint, a.accessKeyId, a.accessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(a.bucketName)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(uri)
	if err != nil {
		return err
	}
	return nil
}

func (a *AliyunOSS) GetObjectFileUrl(uri string) string {
	if strings.Index(a.bucketDomain, "http://") == -1 && strings.Index(a.bucketDomain, "https://") == -1 {
		return "http://" + a.bucketDomain + "/" + uri
	} else {
		return a.bucketDomain + "/" + uri
	}
}
