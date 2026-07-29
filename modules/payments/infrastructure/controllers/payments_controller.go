package paymentscontroller

import (
	"log"
	"net/http"

	paymentsrequestsdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/requests"
	paymentsusecases "github.com/PurpleSavage/monekai-server/modules/payments/application/usecases"
	paymentsmiddlewares "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/middlewares"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)


type PaymentsController struct {
	validator                 *validators.DTOValidator
	authMiddleware            *authmiddlewares.AuthMiddleware
	paymentsMiddleware        *paymentsmiddlewares.PaymentWebhookVerifier
	createPaymentUC           *paymentsusecases.CreatePaymentUseCase
	listCreditPackages        *paymentsusecases.ListCreditPackagesUseCase
	processPaymentWebhookUC   *paymentsusecases.ProcessPaymentWebhookUseCase
	observerBucket            commonports.ObserverBucketPort
}
func NewPaymentsController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	createPaymentUC *paymentsusecases.CreatePaymentUseCase,
	listCreditPackages *paymentsusecases.ListCreditPackagesUseCase,
	paymentsmiddlewares *paymentsmiddlewares.PaymentWebhookVerifier,
	processPaymentWebhookUC *paymentsusecases.ProcessPaymentWebhookUseCase,
	ob commonports.ObserverBucketPort,
) *PaymentsController {
	return &PaymentsController{
		validator:                v,
		authMiddleware:           am,
		createPaymentUC:          createPaymentUC,
		listCreditPackages:       listCreditPackages,
		paymentsMiddleware:       paymentsmiddlewares,
		processPaymentWebhookUC:  processPaymentWebhookUC,
		observerBucket:           ob,
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
	response, err:= sc.createPaymentUC.Execute(
		packageID,
		dtoSession.Id,
		dtoSession.Email,
		r.Context(),
	)
	if err!= nil{
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}


func (sc *PaymentsController) ReceivePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, ok := r.Context().Value(paymentsmiddlewares.WebhookBodyKey).([]byte)
	if !ok || rawBody == nil {
		commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(400, "Bad Request", "Missing webhook body", nil))
		return
	}

	paymentEntity, eventName, err := sc.processPaymentWebhookUC.Execute(r.Context(), rawBody)
	if err != nil {
		log.Printf("ERROR processing payment webhook: %v", err)
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	if paymentEntity != nil {
		data := paymentsrequestsdtos.DataPaymentNotify{
			UserID:    paymentEntity.UserID,
			PaymentID: paymentEntity.ID,
			Credits:   paymentEntity.CreditsPurchased,
			Amount:    paymentEntity.AmountCents,
			Currency:  paymentEntity.Currency,
			Status:    paymentEntity.Status,
		}
		event := commonentities.Event{
			Name: eventName,
			Data: data,
		}
		sc.observerBucket.NotifyObservers(event, "payment_event")
	}

	w.WriteHeader(http.StatusOK)
}

func PaymentsMapRoutes(sc *PaymentsController) chi.Router{
	r := chi.NewRouter()
	r.Use(sc.authMiddleware.AccessToken)
	r.Post("/create", sc.CreatePayment)

	r.Group(func(r chi.Router){
		r.Use(sc.paymentsMiddleware.Verify)
		r.Post("/webhook", sc.ReceivePaymentWebhook)
	})

	return r
}