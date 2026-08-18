package samplerequestsdto

import "gorm.io/datatypes"

type SaveEditSampleDTO struct {
	SampleID string         `json:"sampleId" validate:"required"`
	Effects  datatypes.JSON `json:"effects"`
}
type SaveEditSampleWithURLDTO struct {
	SampleID      string         `json:"sampleId" validate:"required"`
	Effects       datatypes.JSON `json:"effects"`
	FinalAudioURL string         `json:"finalAudioUrl" validate:"required,url"`
}
