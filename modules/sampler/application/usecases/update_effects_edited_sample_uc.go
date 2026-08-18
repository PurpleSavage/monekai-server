package samplerusecases

import (
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
)

type UpdateEffectsEditedSampleUC struct {
	editedRepository samplerports.SamplerEditedPersistencePort
}

func NewUpdateEffectsEditedSampleUC(repo samplerports.SamplerEditedPersistencePort) *UpdateEffectsEditedSampleUC {
	return &UpdateEffectsEditedSampleUC{
		editedRepository: repo,
	}
}

func (uc *UpdateEffectsEditedSampleUC) Execute(id string, dto samplerequestsdto.EffectsDTO) (bool, error) {
	vo, err := commonvalueobjects.CreateEffectsVO(
		dto.Reverb,
		dto.SlowPitch,
		dto.Saturation,
		dto.Delay,
		dto.LowPass,
		dto.HighPass,
		dto.Gain,
		dto.Reverse,
	)
	if err != nil {
		return false, err
	}

	return uc.editedRepository.UpdateEffectsEditedSample(id, *vo)
}
