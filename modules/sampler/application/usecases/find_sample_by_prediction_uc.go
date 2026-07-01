package samplerusecases

import (
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
)
type FindSampleByPredictionUC struct {
	sampleRepository samplerports.SamplerPersistencePort
}

func NewFindSampleByPredictionUC(repo samplerports.SamplerPersistencePort) *FindSampleByPredictionUC {
	return &FindSampleByPredictionUC{
		sampleRepository: repo,
	}
}

func (uc *FindSampleByPredictionUC) Execute(predictionId string) (*samplerentities.SampleEntity,string, error) {
	sample,userId, err := uc.sampleRepository.FindSampleByPredictionId(predictionId)
	if err != nil {
		return nil,"", err
	}
	return sample, userId, nil
}
