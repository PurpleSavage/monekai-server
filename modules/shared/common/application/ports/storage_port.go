package commonports

import (
	"context"
	"io"
)

type StoragePort interface {
	UploadFile(key string, body io.Reader,size int64, contentType string) (string, error)
	GeneratePResignedURL(ctx context.Context,key string, contentType string)(string, error)
}