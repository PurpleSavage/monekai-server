package samplerequestsdto
type ShareSampleRequestDTO struct {
	SampleID string `json:"sampleId" validate:"required"`
}