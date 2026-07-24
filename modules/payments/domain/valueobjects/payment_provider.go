package paymentsvalueobjects

import (
	"strings"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)

type Provider string

const(
	Paddle		Provider = "paddle"
)

type PaymentProvider struct{
	value Provider
}

func CreatePayentProviderVO(
	name string,
)(*PaymentProvider,error){
	nameParse:= strings.TrimSpace(name)
	if nameParse==""{
		return nil, commondomainerrors.NewValidationError(
			"Name",
			"Name is required",
		)
	}
	if Provider(nameParse) != Paddle {
		return nil, commondomainerrors.NewValidationError(
			"Provider payment",
			"Provider payment is not valid",
		)
	}
	return &PaymentProvider{
		value: Provider(nameParse),
	},nil
}

func (p *PaymentProvider) Value() Provider{
	return p.value
}