package accountusecases

import (
	"context"
	"log"

	accountports "github.com/PurpleSavage/monekai-server/modules/account/application/ports"
	accountentities "github.com/PurpleSavage/monekai-server/modules/account/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
)

type ListPaymentHistoryUseCase struct {
	accountRepo accountports.AccountPersistencePort
}

func NewListPaymentHistoryUseCase(
	accountRepo accountports.AccountPersistencePort,
) *ListPaymentHistoryUseCase {
	return &ListPaymentHistoryUseCase{
		accountRepo: accountRepo,
	}
}

func (uc *ListPaymentHistoryUseCase) Execute(
	ctx context.Context,
	userID string,
	page int,
	limit int,
) (*commonentities.PaginatedResult[accountentities.PaymentHistoryEntity], error) {
	var total int
	var payments []accountentities.PaymentHistoryEntity

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		total, err = uc.accountRepo.CountTotalPayments(ctx, userID)
		if err != nil {
			log.Printf("error counting total payments: %v\n", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		payments, err = uc.accountRepo.ListPaymentHistory(ctx, userID, limit, page)
		if err != nil {
			log.Printf("error listing payment history: %v\n", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &commonentities.PaginatedResult[accountentities.PaymentHistoryEntity]{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  payments,
	}, nil
}