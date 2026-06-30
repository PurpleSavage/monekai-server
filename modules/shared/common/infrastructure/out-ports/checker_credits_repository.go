package commonoutports

import (
	"errors"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm" 
)
type CheckerCreditsAdapter struct {
	db *gorm.DB
}
func NewCheckerCreditsAdapter(db *gorm.DB) commonports.CreditsPort{
	return &CheckerCreditsAdapter{
		db: db,
	}
}
func (c *CheckerCreditsAdapter) CheckCredits(email string, creditsNeeded int) (bool,error) {
	var userModel  models.User
	err := c.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {

			return false, globalerrors.NewAppError(404, "User not found", "There is no user with that email", err)
		}
		return false, globalerrors.NewAppError(500, "Error de base de datos", "An error has occurred.", err)
    }
    if userModel.Credits < creditsNeeded{
		return false, nil
	}
	return true, nil	
	
}

func (c *CheckerCreditsAdapter) DecreaseCredits(credits int, email string) (int,error){
	var userModel models.User
	err:= c.db.Model(&models.User{}).
		Where("email = ?",  email).
		Update("credits", userModel.Credits-credits).Error
	if err != nil {
		return 0, err
	}
	return userModel.Credits - credits, nil
}
