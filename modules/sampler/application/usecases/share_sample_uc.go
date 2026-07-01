package samplerusecases

import (
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
)

type ShareSampleUseCase struct {
	audioRepository samplerports.SamplerPersistencePort
}

func NewShareSampleUseCase(
	audioRepository samplerports.SamplerPersistencePort,
) *ShareSampleUseCase {
	return &ShareSampleUseCase{
		audioRepository:audioRepository,
	}
}

func (s *ShareSampleUseCase) Execute(sampleID string,userID string)(*samplerentities.ShareSampleEntity, error) {
	vo,err:=samplervalueobjects.CreateSaveSharedSampleVO(
		sampleID,
		userID,	
		nil,	
		nil,
	)
	if err != nil {
		
		return nil,err
	}
	response,err:= s.audioRepository.SaveSharedSample(*vo)
	if err != nil {
		return nil,err
	}
	return response,nil
}