package communityinfrastructuremappers

import (
	communityresponsesdtos "github.com/PurpleSavage/monekai-server/modules/community/application/dtos/responses"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
)

func BuildListSharedSampleVersionsResponse(versions []communityentities.SharedSampleVersion) []communityresponsesdtos.SharedSampleVersionItemDTO{
	itemsDTO := make([]communityresponsesdtos.SharedSampleVersionItemDTO, len(versions))

	for i, version := range versions {
		itemsDTO[i] = communityresponsesdtos.SharedSampleVersionItemDTO{
			ID:        version.ID,
			Likes:     version.Likes,
			Downloads: version.Downloads,
			CreatedAt: version.CreatedAt,
			SampleVersion: communityresponsesdtos.SampleVersionInfoDTO{
				ID:            version.SampleVersion.ID,
				Effects:       version.SampleVersion.Effects, // Se pasa directo; el motor de JSON usará sus tags internos
				FinalAudioURL: version.SampleVersion.FinalAudioURL,
				SampleName:    version.SampleVersion.SampleName,
				Prompt:        version.SampleVersion.Prompt,
			},
			SharedBy: communityresponsesdtos.SharedByDTO{
				ID:    version.SharedBy.ID,   // Se expone como "userId" en el JSON
				Name:  version.SharedBy.Name,
				Email: version.SharedBy.Email,
			},
		}
	}
	return itemsDTO
}