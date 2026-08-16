package samplerinfrastructuremappers

import (
	"encoding/json"
	"time"

	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplerraws "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/raws"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
)

func ToEditedSampleEntity(
	raw samplerraws.SamplerEditedJoinRaw,
) *samplerentities.EditedSampleEntity {
	effects, err := toEffectsVO(raw.Effects)
	if err != nil {
		effects = commonvalueobjects.EffectsVO{}
	}

	return &samplerentities.EditedSampleEntity{
		ID: raw.ID,
		Sample: samplerentities.SampleEntity{
			Id:           raw.SampleID,
			SampleName:   raw.SampleName,
			Prompt:       raw.Prompt,
			AudioUrl:     &raw.InitialAudioURL,
			Duration:     raw.Duration,
			OutputFormat: raw.OutputFormat,
			ModelVersion: raw.ModelVersion,
			Status:       raw.Status,
			CreatedAt:    raw.CreatedAt.Format(time.RFC3339),
		},
		Effects:       effects,
		FinalAudioURL: raw.FinalAudioURL,
		CreatedAt:     raw.CreatedAt.Format(time.RFC3339),
	}
}

func toEffectsVO(
	effectsJSON []byte,
) (commonvalueobjects.EffectsVO, error) {
	var vo commonvalueobjects.EffectsVO
	if len(effectsJSON) == 0 || string(effectsJSON) == "null" {
		return vo, nil
	}
	if err := json.Unmarshal(effectsJSON, &vo); err != nil {
		return vo, err
	}
	return vo, nil
}
