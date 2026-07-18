package communitycontrollers

import (
	"net/http"

	communityusecases "github.com/PurpleSavage/monekai-server/modules/community/application/usecases"
	communityinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/mappers"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type CommunityController struct {
	validator *validators.DTOValidator
	authmiddleware *authmiddlewares.AuthMiddleware
	listSharedSamples *communityusecases.ListSharedSamplesUC
	listSharedSamplesVersion *communityusecases.ListSharedSamplesVersionUC
	likeToSharedSample *communityusecases.LikeToSharedSampleUC
	downloadToSharedSample *communityusecases.DownloadToSharedSampleUC
}
func NewCommunityController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	listSharedSamples *communityusecases.ListSharedSamplesUC,
	listSharedSamplesVersion *communityusecases.ListSharedSamplesVersionUC,
	likeToSharedSample *communityusecases.LikeToSharedSampleUC,
	downloadToSharedSample *communityusecases.DownloadToSharedSampleUC,
) *CommunityController {
	return &CommunityController{
		validator: v,
		authmiddleware: am,
		listSharedSamples: listSharedSamples,
		listSharedSamplesVersion: listSharedSamplesVersion,
		likeToSharedSample: likeToSharedSample,
		downloadToSharedSample: downloadToSharedSample,
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
	response,error:= nc.listSharedSamples.Execute(r.Context(), paginationVO.Page, paginationVO.Limit)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	responseDTO:=communityinfrastructuremappers.BuildListSharedSamplesResponseDTO(response)
	
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, responseDTO)
}

func (nc *CommunityController) ListSharedEditSamples(w http.ResponseWriter, r *http.Request) {
	page:= r.URL.Query().Get("page")
	limit:= r.URL.Query().Get("limit")
	paginationVO, err := commonvalueobjects.CreatePaginationVO(page, limit)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	response,error:= nc.listSharedSamplesVersion.Execute(r.Context(), paginationVO.Page, paginationVO.Limit)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	responseDTO:=communityinfrastructuremappers.BuildListSharedSampleVersionsResponse(response)
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, responseDTO)
}

func (nc *CommunityController) LikeToSharedSample(w http.ResponseWriter, r *http.Request) {
	sampleID := chi.URLParam(r, "sampleID")
	sampleIDParsed,err:= authvalueobjects.NewUUIDVO(sampleID)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	response, err:= nc.likeToSharedSample.Execute(r.Context(), sampleIDParsed.Value())
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	dtoResponse:=communityinfrastructuremappers.ResponseLikeSharedSampleDTOMapper(response)
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, dtoResponse)
}

func (nc *CommunityController) DownloadSample(w http.ResponseWriter, r *http.Request) {
	sampleID := chi.URLParam(r, "sampleID")
	sampleIDParsed,err:= authvalueobjects.NewUUIDVO(sampleID)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	response, err:= nc.downloadToSharedSample.Execute(r.Context(), sampleIDParsed.Value())
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w,err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}

func CommunityMapRoutes(nc *CommunityController) chi.Router{ 
	r := chi.NewRouter()
	r.Use(nc.authmiddleware.AccessToken)
	r.Get("/samples", nc.ListSharedSamples)
	r.Get("/edit-samples", nc.ListSharedEditSamples)
	r.Patch("/like/{sampleID}",nc.LikeToSharedSample)
	r.Get("/download/{sampleID}", nc.DownloadSample)
	return r
}
