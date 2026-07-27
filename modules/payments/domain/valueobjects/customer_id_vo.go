package paymentsvalueobjects

import (
	"strings"

	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)

type CustomerIDVO struct{
	value string
}
func CreateCustomerIDVO(
	value string,
) (*CustomerIDVO,error){
	stringParsed:= strings.TrimSpace(value)
	if stringParsed==""{
		return nil, commondomainerrors.NewValidationError(
			"Name",
			"Name is required",
		)
	}
	return &CustomerIDVO{
		value:stringParsed,
	}, nil
}
func (c *CustomerIDVO) Value() string{
	return  c.value
}