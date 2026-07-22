package authoutadapters

import (
	"errors"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	authentities "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/entities"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm"
)


type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) authports.UserPersistencePort{
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(data authvalueobjects.CreateUserVO) (*authentities.UserEntity, error) {
    newUser := models.User{
        ExternalID: data.ExternalID,
        Email:      data.Email,
        PhotoURL:   data.PhotoUrl,
        Credits:    100,
    }
    newSession := models.Session{
        RefreshToken: data.RefreshToken,
        UserAgent:    data.UserAgent,
    }

    err := r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&newUser).Error; err != nil {
            return err 
        }
        newSession.UserID = newUser.ID
        if err := tx.Create(&newSession).Error; err != nil {
            return err 
        }

        return nil
    })

    if err != nil {
        return nil, globalerrors.NewAppError(
            500, 
            "Error registering user", 
            "The registration in the database could not be completed", 
            err,
        )
    }
    return &authentities.UserEntity{
        Id:        newUser.ID.String(), 
        Email:     newUser.Email,
        PhotoUrl:  newUser.PhotoURL,
        CreatedAt: newUser.CreatedAt, 
        Credits:   newUser.Credits,
    }, nil
}

func (r *UserRepository) FindUserByEmail(email string)(*authentities.UserEntity,error){
	var userModel  models.User 
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
			// El usuario no existe
			return nil, globalerrors.NewAppError(404, "User not found", "There is no user with that email", err)
		}
		// Error de la bd (conexión, permisos, etc.)
		return nil, globalerrors.NewAppError(500, "Error de base de datos", "Hubo un fallo al conectar con el servidor", err)
    }
	return &authentities.UserEntity{
        Id:        userModel.ID.String(),
        Email:     userModel.Email,
        CreatedAt: userModel.CreatedAt,
        Credits:   userModel.Credits,
        PhotoUrl:  userModel.PhotoURL, 
    }, nil
}

func (r *UserRepository) UpdateSession(token string, userId string) error{
    result := r.db.Model(&models.Session{}).
        Where("user_id = ?", userId).
        Update("refresh_token", token)
    if result.Error != nil {
        return globalerrors.NewAppError(500, "Database Error", "Failed to update session in database", result.Error)
    }
    if result.RowsAffected == 0 {
        // Si no hubo filas afectadas, es porque el user_id no tenía una sesión previa
        return globalerrors.NewAppError(404, "Session Not Found", "No active session found for the given user", nil)
    }
    return  nil
}