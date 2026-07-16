package communityusecases

import (
	"context"
	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
)

type ListSharedSamplesUC struct {
	repo communityports.CommunityPersistencePort
}

func NewListSharedSamplesUC(repo communityports.CommunityPersistencePort) *ListSharedSamplesUC {
	return &ListSharedSamplesUC{repo: repo}
}

func (uc *ListSharedSamplesUC) Execute(
	ctx context.Context, 
	page int, 
	limit int,
) ([]communityentities.SharedSample, error) {
	result, err := uc.repo.ListSharedSamples(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
