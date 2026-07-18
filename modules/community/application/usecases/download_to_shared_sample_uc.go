package communityusecases

import (
	"context"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
	"github.com/google/uuid"
)

type DownloadToSharedSampleUC struct {
	repo communityports.CommunityPersistencePort
}

func NewDownloadToSharedSampleUC(repo communityports.CommunityPersistencePort) *DownloadToSharedSampleUC {
	return &DownloadToSharedSampleUC{
		repo: repo,
	}
}
func (uc *DownloadToSharedSampleUC) Execute(
	ctx context.Context, 
	sampleID uuid.UUID,
) (*communityvalueobjects.DownloadSharedSampleVO, error) {
	resp, err := uc.repo.DownloadToSharedSample(ctx, sampleID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
