package communityusecases

import (
	"context"
	"log"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
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
) (*commonentities.PaginatedResult[communityentities.SharedSampleVersion], error) {
	var total int
	var versions []communityentities.SharedSampleVersion

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		total, err = uc.repo.CountTotalSharedSampleVersions(ctx)
		if err != nil {
			log.Printf("error counting total shared sample versions: %v\n", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		versions, err = uc.repo.ListSharedSamplesVersion(ctx, page, limit)
		if err != nil {
			log.Printf("error listing shared sample versions: %v\n", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &commonentities.PaginatedResult[communityentities.SharedSampleVersion]{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  versions,
	}, nil
}