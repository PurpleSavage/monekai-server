package accountoutadapters

import (
	"context"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	accountports "github.com/PurpleSavage/monekai-server/modules/account/application/ports"
	accountentities "github.com/PurpleSavage/monekai-server/modules/account/domain/entities"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) accountports.AccountPersistencePort {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) CountTotalPayments(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, globalerrors.NewAppError(500, "Database Error", "Error counting total payments", err)
	}
	return int(count), nil
}

func (r *AccountRepository) ListPaymentHistory(
	ctx context.Context,
	userID string,
	limit int,
	page int,
) ([]accountentities.PaymentHistoryEntity, error) {
	var payments []models.Payment

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&payments).Error

	if err != nil {
		return nil, globalerrors.NewAppError(500, "Database Error", "Error listing payment history", err)
	}

	history := make([]accountentities.PaymentHistoryEntity, len(payments))
	for i, payment := range payments {
		history[i] = accountentities.PaymentHistoryEntity{
			ID:          payment.ID.String(),
			CreatedAt:   payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Credits:     payment.CreditsPurchased,
			AmountCents: payment.AmountCents,
			Currency:    payment.Currency,
			Status:      payment.Status,
		}
	}
	return history, nil
}