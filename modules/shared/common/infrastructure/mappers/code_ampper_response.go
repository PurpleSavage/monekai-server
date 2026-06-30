package commoninfrastructuremappers

import (
	"net/http"
	commondomainenums "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/enums"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)

type ErrorConfig struct {
	Status int
	Title  string
}

var domainErrorMap = map[commondomainenums.ErrorCode]ErrorConfig{
	commondomainenums.CodeValidation: {
		Status: http.StatusBadRequest,
		Title:  "Validation Error",
	},
	commondomainenums.CodeConflict: {
		Status: http.StatusConflict,
		Title:  "Conflict",
	},
	commondomainenums.CodeNotFound: {
		Status: http.StatusNotFound,
		Title:  "Not Found",
	},
	commondomainenums.CodeUnauthorized: {
		Status: http.StatusUnauthorized,
		Title:  "Unauthorized",
	},
}
func MapDomainError(err *commondomainerrors.DomainError) *globalerrors.AppError {
 
	config, exists := domainErrorMap[err.Code]

	if !exists {
		return globalerrors.NewAppError(
			http.StatusInternalServerError,
			"Internal Server Error",
			"An unexpected error occurred",
			nil,
		)
	}

	return globalerrors.NewAppError(
		config.Status,
		config.Title,
		err.Message,
		nil,
	)
}