package samplerinfrastructuremappers

import (
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
)

func ToEditedSampleResponseDTO(
	entity *samplerentities.EditedSampleEntity,
) samplerresponsessdtos.EditedSampleResponseDTO {
	return samplerresponsessdtos.EditedSampleResponseDTO{
		ID: entity.ID.String(),
		Sample: samplerresponsessdtos.SampleResponseDTO{
			Id:           entity.Sample.Id.String(),
			SampleName:   entity.Sample.SampleName,
			Prompt:       entity.Sample.Prompt,
			AudioUrl:     entity.Sample.AudioUrl,
			Duration:     entity.Sample.Duration,
			OutputFormat: string(entity.Sample.OutputFormat),
			ModelVersion: string(entity.Sample.ModelVersion),
			Status:       string(entity.Sample.Status),
			CreatedAt:    entity.Sample.CreatedAt,
		},
		Effects: samplerresponsessdtos.EffectsResponseDTO{
			Reverb:     entity.Effects.Reverb,
			SlowPitch:  entity.Effects.SlowPitch,
			Saturation: entity.Effects.Saturation,
			Delay:      entity.Effects.Delay,
			LowPass:    entity.Effects.LowPass,
			HighPass:   entity.Effects.HighPass,
			Gain:       entity.Effects.Gain,
			Reverse:    entity.Effects.Reverse,
		},
		FinalAudioURL: entity.FinalAudioURL,
		CreatedAt:     entity.CreatedAt,
	}
}

func ToEditedSampleResponseDTOList(
	entities []*samplerentities.EditedSampleEntity,
) []samplerresponsessdtos.EditedSampleResponseDTO {
	items := make([]samplerresponsessdtos.EditedSampleResponseDTO, len(entities))
	for i, entity := range entities {
		items[i] = ToEditedSampleResponseDTO(entity)
	}
	return items
}
