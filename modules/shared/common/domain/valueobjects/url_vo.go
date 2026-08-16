package commonvalueobjects

import (
	"net/url"

	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)
type URLVO struct {
	value string
}

func CreateNewURLVO(
	urlValue string,
)(*URLVO,error){
	urlParsed,err := url.ParseRequestURI(urlValue)
	if err!= nil{
		return nil, commondomainerrors.NewValidationError(
			"url",
			"invalid url format",
		)
	}
	if (urlParsed.Scheme != "http" && urlParsed.Scheme != "https") || urlParsed.Host == "" {
		return nil, commondomainerrors.NewValidationError(
			"url",
			"invalid url format, must start with http:// or https://",
		)
	}
	return &URLVO{
		value: urlParsed.String(),
	}, nil
}
func (u *URLVO) Value() string {
	return u.value
}