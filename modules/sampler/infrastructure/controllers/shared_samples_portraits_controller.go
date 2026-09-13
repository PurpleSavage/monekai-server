package samplercontroller

import (
	"net/http"

	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	samplerusecases "github.com/PurpleSavage/monekai-server/modules/sampler/application/usecases"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type SharedSamplesPortraitsController struct {
	authMiddleware       *authmiddlewares.AuthMiddleware
	validator            *validators.DTOValidator
	presignedURLUC       *samplerusecases.GeneratePresignedURLUC
}

func NewSharedSamplesPortraitController(
	am *authmiddlewares.AuthMiddleware,
	v *validators.DTOValidator,
	presignedURLUC *samplerusecases.GeneratePresignedURLUC,
) *SharedSamplesPortraitsController {
	return &SharedSamplesPortraitsController{
		authMiddleware: am,
		validator:      v,
		presignedURLUC: presignedURLUC,
	}
}

func (sp *SharedSamplesPortraitsController) PresignedURL(w http.ResponseWriter, r *http.Request) {
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil))
		return
	}
	dtoSession, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(500, "Internal Error", "Could not parse session data", err))
		return
	}

	dto := samplerequestsdto.PresignedURLRequestDTO{
		Key:  r.URL.Query().Get("key"),
		Name: r.URL.Query().Get("name"),
	}

	if err := sp.validator.ValidateStruct(dto); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	url, key, err := sp.presignedURLUC.Execute(r.Context(), dtoSession.Id, dto.Key, dto.Name)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	response := samplerresponsessdtos.PresignedURLResponseDTO{
		UploadURL: url,
		Key:       key,
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}

func SharedSamplesPortraitsMapRoutes(sp *SharedSamplesPortraitsController, r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(sp.authMiddleware.AccessToken)
		r.Get("/presigned-url", sp.PresignedURL)
	})
}
