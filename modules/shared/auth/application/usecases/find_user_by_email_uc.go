package authusecases

import (
	"strings"

	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	authentities "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/entities"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)

type FindUserByEmailUseCase struct {
	repo authports.UserPersistencePort 
}

func NewFindUserByEmailUseCase(repo authports.UserPersistencePort) *FindUserByEmailUseCase {
	return &FindUserByEmailUseCase{repo: repo}
}

func (uc *FindUserByEmailUseCase) Execute(email string) (*authentities.UserEntity, error) {
	// Limpieza básica de datos
	cleanEmail := strings.TrimSpace(strings.ToLower(email))

	if cleanEmail == "" {
		return nil, globalerrors.NewAppError(400, "Invalid Email", "El correo electrónico es requerido", nil)
	}
	user, err := uc.repo.FindUserByEmail(cleanEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}