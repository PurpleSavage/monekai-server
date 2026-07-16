package communityports

import (
	"context"

	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
)

type CommunityPersistencePort interface {
	ListSharedSamples(
		ctx context.Context,
		page int,
		limit int,
	) ([]communityentities.SharedSample, error)
	
	ListSharedSamplesVersion(
		ctx context.Context,
		page int,
		limit int,
	) ([]communityentities.SharedSampleVersion, error)
}