package samplerinfrastructuremappers

import (
	"time"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	samplerraws "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/raws"
)

func ToResponseUpdateUrlVo(
	result samplerraws.JoinResult,
) *samplervalueobjects.ResponseUpdateurlVo {

	return &samplervalueobjects.ResponseUpdateurlVo{
		UserData: samplervalueobjects.PersonalData{
			Email: result.Email,
			ID:    result.UserID,
		},
		Sample: samplerentities.SampleEntity{
			Id:           result.ID,
			SampleName:   result.SampleName,
			Prompt:       result.Prompt,
			AudioUrl:     result.InitialAudioURL,
			Duration:     result.Duration,
			OutputFormat: result.OutputFormat,
			ModelVersion: result.ModelVersion,
			Status:       result.Status,
			CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		},
	}
}