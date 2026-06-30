package commonmiddlewares

import (
	"net/http"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"gorm.io/gorm"
)

type CheckCreditsMiddleware struct {
	creditService commonports.CreditsPort
	db *gorm.DB
}

func NewCheckCreditsMiddleware(
	creditService commonports.CreditsPort,
) *CheckCreditsMiddleware {
	return &CheckCreditsMiddleware{
		creditService: creditService,
	}
}

func (m *CheckCreditsMiddleware) CheckCredits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		creditsNeeded:= 20
		rawData := r.Context().Value(authmiddlewares.SessionContextKey)
		if rawData == nil {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil))
			return
		}
		dtoSession, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
		if err != nil {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(500, "Internal Error", "Could not parse session data", err))
			return
		}
		hasPermission, err := m.creditService.CheckCredits(dtoSession.Email,creditsNeeded)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !hasPermission {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(403, "Forbidden", "Insufficient credits", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}