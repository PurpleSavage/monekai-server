package communityusecases

import (
	"context"
	"log"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
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
) (*commonentities.PaginatedResult[communityentities.SharedSample], error) {
	var total int
	var samples []communityentities.SharedSample

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		total, err = uc.repo.CountTotalSharedSamples(ctx)
		if err != nil {
			log.Printf("error counting total shared samples: %v\n", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		samples, err = uc.repo.ListSharedSamples(ctx, page, limit)
		if err != nil {
			log.Printf("error listing shared samples: %v\n", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &commonentities.PaginatedResult[communityentities.SharedSample]{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  samples,
	}, nil
}