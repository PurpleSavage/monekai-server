package community

import (
	communitycontrollers "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/controllers"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func CommunityBootstrap(
	db *gorm.DB,
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
) chi.Router{
	controller:= communitycontrollers.NewCommunityController(
		v,
		am,
	)
	return communitycontrollers.CommunityMapRoutes(controller)
}