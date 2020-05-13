package google

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/JerryZhou343/cctool/internal/utils"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type GoogleOSS struct {
	bucketName      string
	credentialsFile string
}

func NewGoogleOSS(bucketName, credentialsFile string) *GoogleOSS {
	return &GoogleOSS{
		bucketName:      bucketName,
		credentialsFile: credentialsFile,
	}

}

func (g *GoogleOSS) UploadFile(srcFilePath string) (uri string, obj string, err error) {

	if !utils.CheckFileExist(srcFilePath) {
		err = status.ErrFileNotExits
		return
	}
	fileName := filepath.Base(srcFilePath)

	ctx := context.Background()
	f, err := os.Open(srcFilePath)
	if err != nil {
		return
	}
	defer f.Close()

	//分日期存储
	date := time.Now()
	year := date.Year()
	month := date.Month()
	day := date.Day()
	obj = strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/" + fileName

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(g.credentialsFile))
	if err != nil {
		return
	}
	wc := client.Bucket(g.bucketName).Object(obj).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return
	}
	if err = wc.Close(); err != nil {
		return
	}

	return g.GetObjectFileUrl(obj), obj, nil
}

func (g *GoogleOSS) GetListBuckets() (ret []string, err error) {
	var (
		attrs *storage.BucketAttrs
	)
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(g.credentialsFile))
	if err != nil {
		return
	}
	it := client.Buckets(ctx, "")

	for {
		attrs, err = it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		ret = append(ret, attrs.Name)
	}
	return
}

func (g *GoogleOSS) DeleteFile(obj string) (err error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(g.credentialsFile))
	if err != nil {
		return
	}
	src := client.Bucket(g.bucketName).Object(obj)

	if err := src.Delete(ctx); err != nil {
		return err
	}

	return
}

func (g *GoogleOSS) GetObjectFileUrl(obj string) string {
	return fmt.Sprintf("gs://%s/%s", g.bucketName, obj)
}
