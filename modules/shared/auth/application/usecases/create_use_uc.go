package authusecases

import (
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	authentities "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/entities"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)
type RegisterUseCase struct{
	repo authports.UserPersistencePort
}
func NewRegisterUseCase(
	repo authports.UserPersistencePort,
) *RegisterUseCase{
	return  &RegisterUseCase{
		repo: repo,
	}
}
func (r*RegisterUseCase) Execute(createUser authvalueobjects.CreateUserVO )(*authentities.UserEntity,error){
	userCreated, err := r.repo.CreateUser(createUser)
    if err != nil {
            return nil, err
    }
	return  userCreated,nil
	
}
