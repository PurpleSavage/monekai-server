package accountports

import (
	"context"

	accountentities "github.com/PurpleSavage/monekai-server/modules/account/domain/entities"
)

type AccountPersistencePort interface {
	CountTotalPayments(ctx context.Context, userID string) (int, error)
	ListPaymentHistory(ctx context.Context, userID string, limit int, page int) ([]accountentities.PaymentHistoryEntity, error)
}