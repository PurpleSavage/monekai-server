package communitycontrollers

import (
	"net/http"

	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type CommunityController struct {
	validator *validators.DTOValidator
	authmiddleware *authmiddlewares.AuthMiddleware
}
func NewCommunityController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
) *CommunityController {
	return &CommunityController{
		validator: v,
		authmiddleware: am,
	}
}

func (nc *CommunityController) ListSharedSamples(w http.ResponseWriter, r *http.Request) {
	
}

func (nc *CommunityController) ListSharedEditSamples(w http.ResponseWriter, r *http.Request) {
	
}

func CommunityMapRoutes(nc *CommunityController) chi.Router{
	r := chi.NewRouter()
	r.Use(nc.authmiddleware.AccessToken)
	r.Get("/samples", nc.ListSharedSamples)
	r.Get("/edit-samples", nc.ListSharedEditSamples)
	return r
}
