package samplerusecases

import (
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
)

type GetEditedSampleByIDUC struct {
	editedRepository samplerports.SamplerEditedPersistencePort
}

func NewGetEditedSampleByIDUC(repo samplerports.SamplerEditedPersistencePort) *GetEditedSampleByIDUC {
	return &GetEditedSampleByIDUC{
		editedRepository: repo,
	}
}

func (uc *GetEditedSampleByIDUC) Execute(id string) (*samplerentities.EditedSampleEntity, error) {
	return uc.editedRepository.GetEditedSampleByID(id)
}
