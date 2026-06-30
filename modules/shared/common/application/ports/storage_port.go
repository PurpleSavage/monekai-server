package commonports

import "io"

type StoragePort interface {
	UploadFile(key string, body io.Reader,size int64, contentType string) (string, error)
}