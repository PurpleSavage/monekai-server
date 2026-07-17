package communityusecases

import (
	"context"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	"github.com/google/uuid"
)

type LikeToSharedSampleUC struct {
	repo communityports.CommunityPersistencePort
}

func NewLikeToSharedSampleUC(repo communityports.CommunityPersistencePort) *LikeToSharedSampleUC {
	return &LikeToSharedSampleUC{repo: repo}
}

func (uc *LikeToSharedSampleUC) Execute(ctx context.Context, sampleID uuid.UUID) (*authvalueobjects.UUIDVO, error) {
	response, err := uc.repo.LikeToSharedSample(ctx, sampleID)
	if err != nil {
		return nil, err
	}
	return response, nil
}