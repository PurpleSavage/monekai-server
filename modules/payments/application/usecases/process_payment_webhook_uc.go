package paymentsusecases

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	paymentsrequestsdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/requests"
	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsentities "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)

type ProcessPaymentWebhookUseCase struct {
	paymentRepo paymentsports.PaymentsPersistencePort
	userRepo    authports.UserPersistencePort
}

func NewProcessPaymentWebhookUseCase(
	paymentRepo paymentsports.PaymentsPersistencePort,
	userRepo authports.UserPersistencePort,
) *ProcessPaymentWebhookUseCase {
	return &ProcessPaymentWebhookUseCase{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

func (uc *ProcessPaymentWebhookUseCase) Execute(ctx context.Context, rawBody []byte) (*paymentsentities.PaymentEntity, commonentities.EventName, error) {
	var payload paymentsrequestsdtos.PaddleWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, "", err
	}

	eventType := payload.EventType
	isPaymentSuccess := eventType == "transaction.billed" || eventType == "transaction.completed"

	if !isPaymentSuccess {
		return nil, commonentities.EventPaymentFailed, nil
	}

	transactionID := payload.Data.ID
	if transactionID == "" {
		return nil, "", nil
	}

	existingPayment, err := uc.paymentRepo.FindPaymentByProviderTransactionID(transactionID)
	if err != nil {
		return nil, "", err
	}
	if existingPayment != nil {
		return existingPayment, commonentities.EventPaymentSuccess, nil
	}

	priceID := ""
	if len(payload.Data.Items) > 0 {
		priceID = payload.Data.Items[0].Price.ID
	}

	creditPackage, err := uc.paymentRepo.GetCreditPackageByPriceID(priceID)
	if err != nil {
		return nil, "", err
	}

	userIDRaw, ok := payload.Data.CustomData["user_id"]
	if !ok {
		return nil, "", nil
	}
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		return nil, "", nil
	}

	user, err := uc.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, "", err
	}

	totalStr := strings.TrimSpace(payload.Data.Total)
	totalFloat, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return nil, "", err
	}
	amountCents := int(math.Round(totalFloat * 100))

	status := strings.ToLower(payload.Data.Status)

	vo, err := paymentsvalueobjects.CreateSavePaymentVO(
		user.Id,
		creditPackage.Id.String(),
		string(paymentsvalueobjects.Paddle),
		transactionID,
		priceID,
		creditPackage.Credits,
		amountCents,
		payload.Data.CurrencyCode,
		status,
		rawBody,
	)
	if err != nil {
		return nil, "", err
	}

	paymentEntity := &paymentsentities.PaymentEntity{
		UserID:                vo.UserID.String(),
		CreditPackageID:       vo.CreditPackageID.String(),
		Provider:              string(vo.Provider),
		ProviderTransactionID: vo.ProviderTransactionID,
		PriceID:               vo.PriceID,
		CreditsPurchased:      vo.CreditsPurchased,
		AmountCents:           vo.AmountCents,
		Currency:              vo.Currency,
		Status:                vo.Status,
		RawPayload:            rawBody,
	}

	if err := uc.paymentRepo.SavePayment(ctx, paymentEntity); err != nil {
		return nil, "", err
	}

	if err := uc.userRepo.AddCredits(user.Id, creditPackage.Credits); err != nil {
		return nil, "", err
	}

	return paymentEntity, commonentities.EventPaymentSuccess, nil
}
