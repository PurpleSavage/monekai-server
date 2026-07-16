package communitycontrollers

import (
	"net/http"

	communityusecases "github.com/PurpleSavage/monekai-server/modules/community/application/usecases"
	communityinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/mappers"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type CommunityController struct {
	validator *validators.DTOValidator
	authmiddleware *authmiddlewares.AuthMiddleware
	listSahredSamples *communityusecases.ListSharedSamplesUC
	listSharedSamplesVersion *communityusecases.ListSharedSamplesVersionUC
}
func NewCommunityController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	listSharedSamples *communityusecases.ListSharedSamplesUC,
	listSharedSamplesVersion *communityusecases.ListSharedSamplesVersionUC,
) *CommunityController {
	return &CommunityController{
		validator: v,
		authmiddleware: am,
		listSahredSamples: listSharedSamples,
		listSharedSamplesVersion: listSharedSamplesVersion,
	}
}

func (nc *CommunityController) ListSharedSamples(w http.ResponseWriter, r *http.Request) {
	page:= r.URL.Query().Get("page")
	limit:= r.URL.Query().Get("limit")
	paginationVO, err := commonvalueobjects.CreatePaginationVO(page, limit)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	response,error:= nc.listSahredSamples.Execute(r.Context(), paginationVO.Page, paginationVO.Limit)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	responseDTO:=communityinfrastructuremappers.BuildListSharedSamplesResponseDTO(response)
	
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, responseDTO)
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
