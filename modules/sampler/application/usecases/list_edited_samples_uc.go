package samplerusecases

import (
	"context"
	"log"

	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
)

type ListEditedSamplesUC struct {
	editedRepository samplerports.SamplerEditedPersistencePort
}

func NewListEditedSamplesUC(repo samplerports.SamplerEditedPersistencePort) *ListEditedSamplesUC {
	return &ListEditedSamplesUC{
		editedRepository: repo,
	}
}

func (uc *ListEditedSamplesUC) Execute(ctx context.Context, userID string, page int, limit int) (*commonentities.PaginatedResult[*samplerentities.EditedSampleEntity], error) {
	var total int
	var samples []*samplerentities.EditedSampleEntity

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		total, err = uc.editedRepository.CountTotalEditedSamples(ctx, userID)
		if err != nil {
			log.Printf("error counting total edited samples: %v\n", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		samples, err = uc.editedRepository.ListSamplesEditedByUsertID(ctx, userID, page, limit)
		if err != nil {
			log.Printf("error listing edited samples: %v\n", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &commonentities.PaginatedResult[*samplerentities.EditedSampleEntity]{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  samples,
	}, nil
}
