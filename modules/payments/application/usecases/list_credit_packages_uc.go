package paymentsusecases

import (
	"context"
	"sync"

	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
)


type ListCreditPackagesUseCase struct{
	paymentsRepo paymentsports.PaymentsPersistencePort
	cacheCreditPackages []*paymentsentites.CreditPackageEntity
	mu                  sync.RWMutex
}

func NewListCreditPackageUseCase(
	paymentsRepo paymentsports.PaymentsPersistencePort,
) *ListCreditPackagesUseCase{
	return  &ListCreditPackagesUseCase{
		paymentsRepo: paymentsRepo,
		cacheCreditPackages: nil,
	}
}

func (l *ListCreditPackagesUseCase) Execute(ctx context.Context)([]*paymentsentites.CreditPackageEntity,error){
	// 1. Bloqueo de LECTURA (varias goroutines pueden leer a la vez)
	l.mu.RLock()
	if len(l.cacheCreditPackages) > 0 {
		cached := l.cacheCreditPackages
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()
	
	// 2. Si no hay caché, vamos a la BD
	listPackages, err := l.paymentsRepo.ListCreditPackages(ctx)
	if err != nil {
		return nil, err
	}
	
	// 3. Bloqueo de ESCRITURA (solo una goroutine escribe)
	l.mu.Lock()
	l.cacheCreditPackages = listPackages
	l.mu.Unlock()
	
	return listPackages, nil
}