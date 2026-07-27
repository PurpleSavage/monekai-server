package paymentsvalueobjects

import (
	"strings"

	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	"github.com/google/uuid"
)



type TransactionVO struct{
	PriceID string
	UserID  uuid.UUID
	CustomerID   string
}
func CreateTransactionVO(
	priceID string,
	userID string,
	CustomerID string,
)(*TransactionVO,error){
	customerParsed:= strings.TrimSpace(CustomerID)
	if customerParsed==""{
		return nil, commondomainerrors.NewValidationError(
			"customerID",
			"CustomerID is required",
		)
	}
	userIDValid,err:=authvalueobjects.NewUUIDVO(userID)
	if err!= nil{
		return nil,err
	}
	return &TransactionVO{
		PriceID: priceID,
		UserID: userIDValid.Value(),
		CustomerID:customerParsed,
	},nil
}
