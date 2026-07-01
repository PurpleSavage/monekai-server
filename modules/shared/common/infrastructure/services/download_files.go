package commonservices

import (
	"io"
	"net/http"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)
func DownloadFile(url string) (io.ReadCloser,int64, error) {
	resp, err := http.Get(url)
	if err != nil {
	 	return nil,0, globalerrors.NewAppError(
			500,
			"Network Error",
			"Could not establish a connection to download the audio file",
			err,
		)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close() 
		               
		return nil,0, globalerrors.NewAppError(502, "External Service Error", "...", nil)
	}
	return resp.Body, resp.ContentLength, nil
}
