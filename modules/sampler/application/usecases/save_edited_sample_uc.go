package samplerusecases

import (
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)

type SaveEditedSampleUC struct {
	editedRepository samplerports.SamplerEditedPersistencePort
}

func NewSaveEditedSampleUC(repo samplerports.SamplerEditedPersistencePort) *SaveEditedSampleUC {
	return &SaveEditedSampleUC{
		editedRepository: repo,
	}
}

func (uc *SaveEditedSampleUC) Execute(dto samplerequestsdto.SaveEditSampleWithURLDTO) (*authvalueobjects.UUIDVO, error) {
	return uc.editedRepository.SaveEditedSample(dto)
}
