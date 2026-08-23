package community

import (
	communityusecases "github.com/PurpleSavage/monekai-server/modules/community/application/usecases"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	communitycontrollers "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/controllers"
	communityoutadapters "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/out-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commoninadapters "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/in-adapters"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func CommunityBootstrap(
	db *gorm.DB,
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
) chi.Router {
	communityRepo := communityoutadapters.NewCommunityRepository(db)
	//use cases
	listSharedSamples := communityusecases.NewListSharedSamplesUC(communityRepo)
	listSharedSamplesVersion := communityusecases.NewListSharedSamplesVersionUC(communityRepo)
	likeToSharedSample := communityusecases.NewLikeToSharedSampleUC(communityRepo)
	downloadToSharedSample := communityusecases.NewDownloadToSharedSampleUC(communityRepo)
	latestSharedSamplesCache := commoninadapters.NewLocalCache[[]communityentities.SharedSample]()
	latestSharedEffectSamplesCache := commoninadapters.NewLocalCache[[]communityentities.SharedSampleVersion]()
	getLatestSharedSamples := communityusecases.NewGetLatestSharedSamplesUC(communityRepo, latestSharedSamplesCache)
	getLatestSharedEffectSamples := communityusecases.NewGetLatestSharedEffectSamplesUC(communityRepo, latestSharedEffectSamplesCache)

	controller := communitycontrollers.NewCommunityController(
		v,
		am,
		listSharedSamples,
		listSharedSamplesVersion,
		likeToSharedSample,
		downloadToSharedSample,
		getLatestSharedSamples,
		getLatestSharedEffectSamples,
	)
	return communitycontrollers.CommunityMapRoutes(controller)
}
