package communityusecases

import (
	"context"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
)
type ListSharedSamplesVersionUC struct {
	repo communityports.CommunityPersistencePort
}
func NewListSharedSamplesVersionUC(repo communityports.CommunityPersistencePort) *ListSharedSamplesVersionUC {
	return &ListSharedSamplesVersionUC{repo: repo}
}
func (uc *ListSharedSamplesVersionUC) Execute(
	ctx context.Context,
	page int,
	limit int,
) ([]communityentities.SharedSampleVersion, error) {
	respons, err := uc.repo.ListSharedSamplesVersion(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return respons, nil
}
