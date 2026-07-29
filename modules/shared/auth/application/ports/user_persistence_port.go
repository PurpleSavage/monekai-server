package authports

import (
	authentities "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/entities"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)

type UserPersistencePort interface {
	CreateUser(data authvalueobjects.CreateUserVO)(*authentities.UserEntity,error)
	FindUserByEmail(email string)(*authentities.UserEntity,error)
	FindUserByID(userID string) (*authentities.UserEntity, error)
	UpdateSession(token string, userId string) error
	UpdateCustomerID(customerID string,userID string)(string,error)
	AddCredits(userID string, credits int) error
}