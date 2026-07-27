package paymentscontroller

import (
	"net/http"

	paymentsusecases "github.com/PurpleSavage/monekai-server/modules/payments/application/usecases"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)


type PaymentsController struct {
	validator  *validators.DTOValidator
	authMiddleware *authmiddlewares.AuthMiddleware
	createPaymentUC *paymentsusecases.CreatePaymentUseCase
	listCreditPackages *paymentsusecases.ListCreditPackagesUseCase
}
func NewPaymentsController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	createPaymentUC *paymentsusecases.CreatePaymentUseCase,
	listCreditPackages *paymentsusecases.ListCreditPackagesUseCase,
) *PaymentsController {
	return &PaymentsController{
		validator:v,
		authMiddleware: am,
		createPaymentUC: createPaymentUC,
		listCreditPackages: listCreditPackages,
	}
}
func (sc *PaymentsController) CreatePayment(w http.ResponseWriter, r *http.Request) {
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
	packageID := r.URL.Query().Get("packageId")
	response, err:= sc.createPaymentUC.Execute(packageID,r.Context())
	if err!= nil{
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	//TODO: modificar
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}


func (sc *PaymentsController) ReceivePaymentWebhook(w http.ResponseWriter, r *http.Request){
	
}

func PaymentsMapRoutes(sc *PaymentsController) chi.Router{
	r := chi.NewRouter()
	r.Use(sc.authMiddleware.AccessToken)
	r.Post("/create", sc.CreatePayment)
	r.Post("/webhook", sc.ReceivePaymentWebhook)
	return r
}