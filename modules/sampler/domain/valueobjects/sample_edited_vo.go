package samplervalueobjects

import (
	"encoding/json"
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)


type SampleEditedVO struct {
	
	SampleID uuid.UUID 
	Effects  commonvalueobjects.EffectsVO
	FinalAudioURL string 
}

func CreateSampleEditedVO(

	sampleID string,
	effectsJSON datatypes.JSON,
	audioUrl string,
)(*SampleEditedVO,error){

	sampleIDParsed, err:= authvalueobjects.NewUUIDVO(sampleID)
	if err!=nil {
		return nil ,err
	}
	audioURL,err := commonvalueobjects.CreateNewURLVO(audioUrl)
	if err!=nil {
		return nil ,err
	}

	var dto samplerequestsdto.EffectsDTO
	if len(effectsJSON) > 0 && string(effectsJSON) != "null" {
		if err := json.Unmarshal(effectsJSON, &dto); err != nil {
			return nil, commondomainerrors.NewValidationError(
				"effects",
				"invalid json structure for effects",
			)
		}
	}

	effectsVO, err := commonvalueobjects.CreateEffectsVO(
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
		return nil, err
	}
	
	return &SampleEditedVO{
		SampleID:      sampleIDParsed.Value(),
		Effects:       *effectsVO,
		FinalAudioURL: audioURL.Value(),
	}, nil
}

