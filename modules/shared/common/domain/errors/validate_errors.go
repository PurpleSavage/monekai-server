package commondomainerrors

import commondomainenums "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/enums"

type DomainError struct {
	Code    commondomainenums.ErrorCode
	Field   string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewValidationError(field, message string) error {
	return &DomainError{
		Code:    commondomainenums.CodeValidation,
		Field:   field,
		Message: message,
	}
}