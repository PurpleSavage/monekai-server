package paymentsvalueobjects

import (
	"strings"

	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	"github.com/google/uuid"
)

type SavePaymentVO struct {
	UserID                uuid.UUID
	CreditPackageID       uuid.UUID
	Provider              Provider
	ProviderTransactionID string
	PriceID               string
	CreditsPurchased      int
	AmountCents           int
	Currency              string
	Status                string
	RawPayload            []byte
}

func CreateSavePaymentVO(
	userID string,
	creditPackageID string,
	provider string,
	providerTransactionID string,
	priceID string,
	creditsPurchased int,
	amountCents int,
	currency string,
	status string,
	rawPayload []byte,
) (*SavePaymentVO, error) {
	userUUID, err := authvalueobjects.NewUUIDVO(userID)
	if err != nil {
		return nil, err
	}

	packageUUID, err := authvalueobjects.NewUUIDVO(creditPackageID)
	if err != nil {
		return nil, err
	}

	providerParsed := strings.TrimSpace(provider)
	if providerParsed == "" {
		return nil, commondomainerrors.NewValidationError(
			"provider",
			"provider is required",
		)
	}
	providerVO, err := CreatePayentProviderVO(provider)
	if err != nil {
		return nil, err
	}

	transactionID := strings.TrimSpace(providerTransactionID)
	if transactionID == "" {
		return nil, commondomainerrors.NewValidationError(
			"providerTransactionID",
			"provider transaction ID is required",
		)
	}

	priceID = strings.TrimSpace(priceID)
	if priceID == "" {
		return nil, commondomainerrors.NewValidationError(
			"priceID",
			"price ID is required",
		)
	}

	if creditsPurchased <= 0 {
		return nil, commondomainerrors.NewValidationError(
			"creditsPurchased",
			"credits must be greater than zero",
		)
	}

	if amountCents <= 0 {
		return nil, commondomainerrors.NewValidationError(
			"amountCents",
			"amount must be greater than zero",
		)
	}

	currency = strings.TrimSpace(strings.ToUpper(currency))
	if currency == "" {
		return nil, commondomainerrors.NewValidationError(
			"currency",
			"currency is required",
		)
	}

	status = strings.TrimSpace(strings.ToLower(status))
	if status == "" {
		return nil, commondomainerrors.NewValidationError(
			"status",
			"status is required",
		)
	}

	return &SavePaymentVO{
		UserID:                userUUID.Value(),
		CreditPackageID:       packageUUID.Value(),
		Provider:              providerVO.Value(),
		ProviderTransactionID: transactionID,
		PriceID:               priceID,
		CreditsPurchased:      creditsPurchased,
		AmountCents:           amountCents,
		Currency:              currency,
		Status:                status,
		RawPayload:            rawPayload,
	}, nil
}
