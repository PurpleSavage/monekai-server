package authvalueobjects

import (
	"strings"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)


type TokenVO struct {
	value string
}

func NewTokenVO(value string) (*TokenVO, error) {
	if strings.TrimSpace(value) == "" {

		return nil, commondomainerrors.NewValidationError(
			"token",
			"token cannot be empty",
		)
	}

	return &TokenVO{
		value: value,
	}, nil
}

func (t TokenVO) Value() string {
	return t.value
}