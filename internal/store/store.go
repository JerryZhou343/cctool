package store

type Store interface {
	UploadFile(localFilePath string) (url string, objName string, err error)
	DeleteFile(uri string) error
}
