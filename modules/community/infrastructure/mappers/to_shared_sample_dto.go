package communityinfrastructuremappers

import (
	communityresponsesdtos "github.com/PurpleSavage/monekai-server/modules/community/application/dtos/responses"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
)
func BuildListSharedSamplesResponseDTO(samples []communityentities.SharedSample) []communityresponsesdtos.SharedSampleItemDTO {
	itemsDTO := make([]communityresponsesdtos.SharedSampleItemDTO, len(samples))

	for i, sample := range samples {
		itemsDTO[i] = communityresponsesdtos.SharedSampleItemDTO{
			ID:        sample.ID,
			Likes:     sample.Likes,
			Downloads: sample.Downloads,
			CreatedAt: sample.CreatedAt,
			Sample: communityresponsesdtos.SampleInfoDTO{
				ID:              sample.Sample.ID,
				SampleName:      sample.Sample.SampleName,
				InitialAudioURL: sample.Sample.InitialAudioURL,
				Prompt:          sample.Sample.Prompt,
				Duration:        sample.Sample.Duration,
			},
			SharedBy: communityresponsesdtos.SharedByDTO{
				ID:    sample.SharedBy.ID, // Se convertirá en "userId" en el JSON
				Name:  sample.SharedBy.Name,
				Email: sample.SharedBy.Email,
			},
		}
	}

	return itemsDTO
}