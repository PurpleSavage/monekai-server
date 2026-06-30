package commoninfrastructuremappers

import (
	"encoding/json"
	"net/http"

	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, err error) {

	// Error HTTP
	if appErr, ok := err.(*globalerrors.AppError); ok {
		RespondWithJSON(w, appErr.Status, map[string]any{
			"title":   appErr.Title,
			"message": appErr.Message,
			"status":  appErr.Status,
		})
		return
	}

	// Error de dominio
	if domainErr, ok := err.(*commondomainerrors.DomainError); ok {

		appErr := MapDomainError(domainErr)

		RespondWithJSON(w, appErr.Status, map[string]any{
			"title":   appErr.Title,
			"message": appErr.Message,
			"status":  appErr.Status,
			"field":   domainErr.Field, // opcional
		})

		return
	}

	// Error inesperado
	RespondWithJSON(w, http.StatusInternalServerError, map[string]any{
		"title":   "Internal Server Error",
		"message": "An unexpected error occurred",
		"status":  http.StatusInternalServerError,
	})
}