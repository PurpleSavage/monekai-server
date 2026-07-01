package samplerusecases

import (
	"context"
	"log"

	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)
type ListSampleUseCase struct { 
	audioRepository samplerports.SamplerPersistencePort
}

func NewListSampleUseCase(
	audioRepository samplerports.SamplerPersistencePort,
) *ListSampleUseCase {
	return &ListSampleUseCase{
		audioRepository:audioRepository,
	}
}

func (l *ListSampleUseCase) Execute(
	ctx context.Context,
	userID string,
	page int,
	limit int,
)(*commonentities.PaginatedResult[*samplerentities.SampleEntity],error){

	total, err:= l.audioRepository.CountTotalSamples(ctx,userID)

	if err != nil {
		log.Printf("error counting total samples: %v\n", err)
		return nil, err
	}

	samples,err := l.audioRepository.ListSamples(ctx,userID,page,limit)

	if err!=nil{
		log.Printf("error listing samples: %v\n", err)
		return nil,err
	}
	
	return &commonentities.PaginatedResult[*samplerentities.SampleEntity]{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  samples,
	}, nil
}