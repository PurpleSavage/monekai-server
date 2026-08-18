package samplerusecases

import (
	"context"
	"log"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
)

type ListSampleUseCase struct {
	audioRepository samplerports.SamplerPersistencePort
}

func NewListSampleUseCase(
	audioRepository samplerports.SamplerPersistencePort,
) *ListSampleUseCase {
	return &ListSampleUseCase{
		audioRepository: audioRepository,
	}
}

func (l *ListSampleUseCase) Execute(
	ctx context.Context,
	userID string,
	page int,
	limit int,
) (*commonentities.PaginatedResult[*samplerentities.SampleEntity], error) {
	var totalNotifications int
	var samples []*samplerentities.SampleEntity
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func()error{
		var err error
		totalNotifications, err = l.audioRepository.CountTotalSamples(ctx, userID)
		return err
	})
	g.Go(func() error{
		var err error
		samples, err= l.audioRepository.ListSamples(ctx, userID, page, limit)
		return err
	})
	if err := g.Wait(); err != nil {
		log.Printf("error counting total samples: %v\n", err)
		return nil, err
	}


	return &commonentities.PaginatedResult[*samplerentities.SampleEntity]{
		Total: totalNotifications,
		Limit: limit,
		Page:  page,
		Data:  samples,
	}, nil
}
