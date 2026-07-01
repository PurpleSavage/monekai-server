package samplercontroller

import (
	samplerusecases "github.com/PurpleSavage/monekai-server/modules/sampler/application/usecases"
	samplermiddlewares "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/middlewares"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type SamplerController struct {
	creditsMiddleware *commonmiddlewares.CheckCreditsMiddleware
	authMiddleware *authmiddlewares.AuthMiddleware
	replicateMiddleware *samplermiddlewares.ReplicateMiddlewareWebhook
	generateSongUseCase *samplerusecases.GenerateSongUseCase
	updateUrlSampleGeneratedUseCase *samplerusecases.UpdateUrlSampleGenerated
	validator           *validators.DTOValidator
	sseManager *audiosse.SSEManager
	listSampleUseCase *samplerusecases.ListSampleUseCase
	shareSampleUseCase *samplerusecases.ShareSampleUseCase
	saveNotificationUseCase *samplerusecases.SaveNotificationUC
	findSampleByPredictionUseCase *samplerusecases.FindSampleByPredictionUC
}
func NewSamplerController(
	cm *privatemiddlewares.CheckCreditsMiddleware,
	am *privatemiddlewares.AuthMiddleware,
	gs *audiousecases.GenerateSongUseCase, 
	v *sharedvalidators.DTOValidator,
	sse *audiosse.SSEManager,
	rep *privatemiddlewares.ReplicateMiddlewareWebhook,
	usg *audiousecases.UpdateUrlSampleGenerated, 
	ls *audiousecases.ListSampleUseCase,
	ss *audiousecases.ShareSampleUseCase,
	sn *audiousecases.SaveNotificationUC,
	fs *audiousecases.FindSampleByPredictionUC,
) *SamplerController {
	return &SamplerController{
		creditsMiddleware: cm,
		authMiddleware: am,
		generateSongUseCase:gs,
		validator:           v,
		sseManager:sse,
		replicateMiddleware: rep,
		updateUrlSampleGeneratedUseCase:usg,
		listSampleUseCase:ls,
		shareSampleUseCase:ss,
		saveNotificationUseCase:sn,
		findSampleByPredictionUseCase:fs,
	}
}

func SamplerMapRoutes(controller *SamplerController) chi.Router{
	r := chi.NewRouter()
	
	return r
}
