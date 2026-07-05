package samplercontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	samplerusecases "github.com/PurpleSavage/monekai-server/modules/sampler/application/usecases"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	samplermiddlewares "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/middlewares"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonresponsesdtos "github.com/PurpleSavage/monekai-server/modules/shared/common/application/dtos/responses"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	commonmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/middlewares"
	commonutils "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/utils"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type SamplerController struct {
	creditsMiddleware *commonmiddlewares.CheckCreditsMiddleware
	authMiddleware *authmiddlewares.AuthMiddleware
	replicateMiddleware *samplermiddlewares.ReplicateMiddlewareWebhook
	generateSampleUseCase *samplerusecases.GenerateSampleUseCase
	updateUrlSampleGeneratedUseCase *samplerusecases.UpdateUrlSampleGenerated
	validator           *validators.DTOValidator
	listSampleUseCase *samplerusecases.ListSampleUseCase
	shareSampleUseCase *samplerusecases.ShareSampleUseCase
	findSampleByPredictionUseCase *samplerusecases.FindSampleByPredictionUC
	observerBucket commonports.ObserverBucketPort
}
func NewSamplerController(
	cm *commonmiddlewares.CheckCreditsMiddleware,
	am *authmiddlewares.AuthMiddleware,
	gs *samplerusecases.GenerateSampleUseCase, 
	v *validators.DTOValidator,
	rep *samplermiddlewares.ReplicateMiddlewareWebhook,
	usg *samplerusecases.UpdateUrlSampleGenerated, 
	ls *samplerusecases.ListSampleUseCase,
	ss *samplerusecases.ShareSampleUseCase,
	fs *samplerusecases.FindSampleByPredictionUC,
	ob commonports.ObserverBucketPort,
) *SamplerController {
	return &SamplerController{
		creditsMiddleware: cm,
		authMiddleware: am,
		generateSampleUseCase:gs,
		validator:           v,
		replicateMiddleware: rep,
		updateUrlSampleGeneratedUseCase:usg,
		listSampleUseCase:ls,
		shareSampleUseCase:ss,
		findSampleByPredictionUseCase:fs,
		observerBucket: ob,
	}
}
// controlador para crear un sample: retorna una respuesta automáticamente mientras se porcesa la solioitud 
// en replicate
func (h *SamplerController) CreateSong(w http.ResponseWriter, r *http.Request){
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
	var dto samplerequestsdto.GenerateSampleDTO
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
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
	if err := h.validator.ValidateStruct(dto); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	response, err := h.generateSampleUseCase.Execute(dto,dtoSession.Id)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)

}




func (h *SamplerController) ReceiveWebhook(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload samplerequestsdto.SongWebhookResponse
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("DECODE ERROR: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	} 
	defer r.Body.Close()
	if payload.Status != string(samplerenums.Succeeded) {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ==========================================
	// CASO: REPLICATE DEVOLVIÓ UN ERROR
	// ==========================================
	
	
}




func ( h*SamplerController) ListSamples(w http.ResponseWriter, r *http.Request){
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		commoninfrastructuremappers.RespondWithError(
			w, 
			globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil),
		)
		return
	}

	page:= r.URL.Query().Get("page")
	pageNumber,err:= strconv.Atoi(page)
	if err!=nil{
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusBadRequest,
				"Bad Request",
				"The 'page' query parameter must be a valid positive integer",
				err,
			),
		)
		return
	}
	limit := r.URL.Query().Get("limit")

	limitNumber,err:= strconv.Atoi(limit)
	if err != nil{
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(
				http.StatusBadRequest,
				"Bad Request",
				"The 'page' query parameter must be a valid positive integer",
				err,
			),
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
	result, err := h.listSampleUseCase.Execute(r.Context(), dtoSession.Id,pageNumber, limitNumber)

	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	sampleDTOs := make([]samplerresponsessdtos.SampleResponseDTO, len(result.Data))
	for i, entity := range result.Data {
		sampleDTOs[i] = samplerresponsessdtos.SampleResponseDTO{
			Id: entity.Id.String(),
			SampleName:entity.SampleName,
			Prompt:entity.Prompt,
			AudioUrl:entity.AudioUrl,
			Duration:entity.Duration,
			OutputFormat:string(entity.OutputFormat),
			ModelVersion:string(entity.ModelVersion),
			Status:string(entity.Status),
			CreatedAt:entity.CreatedAt,
		}
	}

	response := commonresponsesdtos.PaginatedResponse[samplerresponsessdtos.SampleResponseDTO]{
		Total: result.Total,
		Limit: result.Limit,
		Page:  result.Page,
		Data:  sampleDTOs,
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}




func (h *SamplerController) ShareSample(w http.ResponseWriter, r *http.Request){
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
	var dto samplerequestsdto.ShareSampleRequestDTO
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
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

	if err := h.validator.ValidateStruct(dto); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	entity, err := h.shareSampleUseCase.Execute(dto.SampleID,dtoSession.Id)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	response := samplerresponsessdtos.ShareSampleResponseDTO{
		Id:              entity.Id.String(),
		SampleId:        commonutils.UUIDPtrToStringPtr(entity.SampleId),
		SampleVersionId: commonutils.UUIDPtrToStringPtr(entity.SampleVersionId),
		UserId:          entity.UserId,
		Likes:           entity.Likes,
		Downloads:       entity.Downloads,
		CreatedAt:       entity.CreatedAt,
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}
func (h *SamplerController) SaveSampleVersion(w http.ResponseWriter, r *http.Request){
	
}

func SamplerMapRoutes(sc *SamplerController) chi.Router{
	r := chi.NewRouter()

	r.Use(sc.authMiddleware.AccessToken)

    r.Group(func(r chi.Router) {
        r.Use(sc.creditsMiddleware.CheckCredits)
        r.Post("/create", sc.CreateSong)
    })

    r.Get("/samples", sc.ListSamples)
    r.Post("/share-sample", sc.ShareSample)
    r.Post("/sample-version", sc.SaveSampleVersion)

    r.Group(func(r chi.Router) {
        r.Use(sc.replicateMiddleware.VeriFyWebhook)
        r.Post("/webhook/songs", sc.ReceiveWebhook)
    })
    
	return r
}
