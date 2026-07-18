package communityports

import (
	"context"

	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	"github.com/google/uuid"
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

	LikeToSharedSample(
		ctx context.Context,
		sampleID uuid.UUID,
	)(*authvalueobjects.UUIDVO,error)

	DownloadToSharedSample(
		ctx context.Context,
		sampleID uuid.UUID,
	) (*communityvalueobjects.DownloadSharedSampleVO, error)
}