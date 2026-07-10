package sampler

import (
	samplerusecases "github.com/PurpleSavage/monekai-server/modules/sampler/application/usecases"
	samplercontroller "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/controllers"
	samplermiddlewares "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/middlewares"
	sampleroutadapters "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/out-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/middlewares"
	commonoutadapters "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/out-adapters"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)
func SamplerBootstrap(
	db *gorm.DB,
	ob commonports.ObserverBucketPort,
	v *validators.DTOValidator,
 	authmiddleware *authmiddlewares.AuthMiddleware,
) chi.Router{

	//adapterts
	storageService:=commonoutadapters.NewCloudFlareAdapterService()
	songService,_:=sampleroutadapters.NewReplicateAdapterService()
	samplerRepo:=sampleroutadapters.NewSamplerrepository(db)
	checkerCreditsService:= commonoutadapters.NewCheckerCreditsAdapter(db)

	// use cases
	generateSampleUC:=samplerusecases.NewGeneratorSampleUseCase(songService,samplerRepo)
	updateUrlSampleUC:=samplerusecases.NewUpdateUrlSampleGenerated(samplerRepo,storageService)
	listSamplesUC:= samplerusecases.NewListSampleUseCase(samplerRepo)
	shareSampleUC:=samplerusecases.NewShareSampleUseCase(samplerRepo)
	findSampleByPredictionUC:=samplerusecases.NewFindSampleByPredictionUC(samplerRepo)

	// middlewaress
	// middlewares 
	replicateMiddleware := samplermiddlewares.NewReplicateMiddlewareWebhook()
	creditsMiddleware := commonmiddlewares.NewCheckCreditsMiddleware(checkerCreditsService)

	
	controller := samplercontroller.NewSamplerController(
		creditsMiddleware,
		authmiddleware,
		generateSampleUC,
		v,
		replicateMiddleware,
		updateUrlSampleUC,
		listSamplesUC,
		shareSampleUC,
		findSampleByPredictionUC,
		ob,
	)
	return samplercontroller.SamplerMapRoutes(controller)
}