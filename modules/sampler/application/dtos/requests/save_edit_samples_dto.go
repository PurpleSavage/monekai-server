package samplerequestsdto

import "gorm.io/datatypes"

type SaveEditSampleDTO struct{
	SampleID string
	Effects datatypes.JSON
}
type SaveEditSampleWithURLDTO struct{
	SampleID string
	Effects datatypes.JSON
	FinalAudioURL string 
}