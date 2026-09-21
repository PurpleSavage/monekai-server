package accountcontrollers

import (
	"net/http"

	accountresponsesdtos "github.com/PurpleSavage/monekai-server/modules/account/application/dtos/responses"
	accountusecases "github.com/PurpleSavage/monekai-server/modules/account/application/usecases"
	accountinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/account/infrastructure/mappers"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonresponsesdtos "github.com/PurpleSavage/monekai-server/modules/shared/common/application/dtos/responses"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/go-chi/chi/v5"
)

type AccountController struct {
	authMiddleware         *authmiddlewares.AuthMiddleware
	listPaymentHistory     *accountusecases.ListPaymentHistoryUseCase
}

func NewAccountController(
	am *authmiddlewares.AuthMiddleware,
	listPaymentHistory *accountusecases.ListPaymentHistoryUseCase,
) *AccountController {
	return &AccountController{
		authMiddleware:     am,
		listPaymentHistory: listPaymentHistory,
	}
}

func (c *AccountController) ListPaymentHistory(w http.ResponseWriter, r *http.Request) {
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

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	paginationVO, err := commonvalueobjects.CreatePaginationVO(page, limit)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	response, err := c.listPaymentHistory.Execute(r.Context(), dtoSession.Id, paginationVO.Page, paginationVO.Limit)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	responseDTO := commonresponsesdtos.PaginatedResponse[accountresponsesdtos.PaymentHistoryItemDTO]{
		Total: response.Total,
		Limit: response.Limit,
		Page:  response.Page,
		Data:  accountinfrastructuremappers.BuildListPaymentHistoryResponseDTO(response.Data),
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, responseDTO)
}

func AccountMapRoutes(ac *AccountController) chi.Router {
	r := chi.NewRouter()
	r.Use(ac.authMiddleware.AccessToken)
	r.Get("/payments/history", ac.ListPaymentHistory)
	return r
}