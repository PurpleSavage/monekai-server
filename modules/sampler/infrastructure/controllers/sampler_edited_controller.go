package samplercontroller

import (
	"encoding/json"
	"net/http"

	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	samplerusecases "github.com/PurpleSavage/monekai-server/modules/sampler/application/usecases"
	samplerinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/mappers"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonresponsesdtos "github.com/PurpleSavage/monekai-server/modules/shared/common/application/dtos/responses"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type SamplerEditedController struct {
	validator                   *validators.DTOValidator
	authMiddleware              *authmiddlewares.AuthMiddleware
	saveEditedSampleUC          *samplerusecases.SaveEditedSampleUC
	getEditedSampleByIDUC       *samplerusecases.GetEditedSampleByIDUC
	updateURLEditedSampleUC     *samplerusecases.UpdateURLEditedSampleUC
	updateEffectsEditedSampleUC *samplerusecases.UpdateEffectsEditedSampleUC
	listEditedSamplesUC         *samplerusecases.ListEditedSamplesUC
}

func NewSamplerEditedController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	save *samplerusecases.SaveEditedSampleUC,
	get *samplerusecases.GetEditedSampleByIDUC,
	updateURL *samplerusecases.UpdateURLEditedSampleUC,
	updateEffects *samplerusecases.UpdateEffectsEditedSampleUC,
	list *samplerusecases.ListEditedSamplesUC,
) *SamplerEditedController {
	return &SamplerEditedController{
		validator:                   v,
		authMiddleware:              am,
		saveEditedSampleUC:          save,
		getEditedSampleByIDUC:       get,
		updateURLEditedSampleUC:     updateURL,
		updateEffectsEditedSampleUC: updateEffects,
		listEditedSamplesUC:         list,
	}
}

func (h *SamplerEditedController) SaveEditedSample(w http.ResponseWriter, r *http.Request) {
	var dto samplerequestsdto.SaveEditSampleWithURLDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusBadRequest,
				"Bad Request",
				"Malformed JSON payload",
				err,
			),
		)
		return
	}
	defer r.Body.Close()
	if err := h.validator.ValidateStruct(dto); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	id, err := h.saveEditedSampleUC.Execute(dto)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

func (h *SamplerEditedController) GetEditedSampleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entity, err := h.getEditedSampleByIDUC.Execute(id)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	response := samplerinfrastructuremappers.ToEditedSampleResponseDTO(entity)
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}

func (h *SamplerEditedController) UpdateURLEditedSample(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var dto samplerequestsdto.UpdateEditedSampleURLDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusBadRequest,
				"Bad Request",
				"Malformed JSON payload",
				err,
			),
		)
		return
	}
	defer r.Body.Close()
	if err := h.validator.ValidateStruct(dto); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	updated, err := h.updateURLEditedSampleUC.Execute(id, dto.FinalAudioURL)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, map[string]bool{"updated": updated})
}

func (h *SamplerEditedController) UpdateEffectsEditedSample(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var dto samplerequestsdto.EffectsDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusBadRequest,
				"Bad Request",
				"Malformed JSON payload",
				err,
			),
		)
		return
	}
	defer r.Body.Close()

	updated, err := h.updateEffectsEditedSampleUC.Execute(id, dto)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, map[string]bool{"updated": updated})
}

func (h *SamplerEditedController) ListEditedSamples(w http.ResponseWriter, r *http.Request) {
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil),
		)
		return
	}
	dtoSession, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusUnauthorized,
				"Unauthorized",
				"Invalid session data",
				err,
			),
		)
		return
	}

	paginationVO, err := commonvalueobjects.CreatePaginationVO(
		r.URL.Query().Get("page"),
		r.URL.Query().Get("limit"),
	)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	result, err := h.listEditedSamplesUC.Execute(
		r.Context(),
		dtoSession.Id,
		paginationVO.Page,
		paginationVO.Limit,
	)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	response := commonresponsesdtos.PaginatedResponse[samplerresponsessdtos.EditedSampleResponseDTO]{
		Total: result.Total,
		Limit: result.Limit,
		Page:  result.Page,
		Data:  samplerinfrastructuremappers.ToEditedSampleResponseDTOList(result.Data),
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}

func SamplerEditedMapRoutes(se *SamplerEditedController) chi.Router {
	r := chi.NewRouter()
	r.Use(se.authMiddleware.AccessToken)

	r.Post("/edit-samples", se.SaveEditedSample)
	r.Get("/edit-samples", se.ListEditedSamples)
	r.Get("/edit-samples/{id}", se.GetEditedSampleByID)
	r.Patch("/edit-samples/{id}/url", se.UpdateURLEditedSample)
	r.Patch("/edit-samples/{id}/effects", se.UpdateEffectsEditedSample)

	return r
}
