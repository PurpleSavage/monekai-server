package samplerusecases

import (
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
)

type UpdateURLEditedSampleUC struct {
	editedRepository samplerports.SamplerEditedPersistencePort
}

func NewUpdateURLEditedSampleUC(repo samplerports.SamplerEditedPersistencePort) *UpdateURLEditedSampleUC {
	return &UpdateURLEditedSampleUC{
		editedRepository: repo,
	}
}

func (uc *UpdateURLEditedSampleUC) Execute(id string, url string) (bool, error) {
	return uc.editedRepository.UpdateURLEditedSample(id, url)
}
