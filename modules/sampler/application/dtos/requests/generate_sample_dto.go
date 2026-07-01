package samplerequestsdto

import samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
type GenerateSampleDTO struct {

	// User prompt
	Prompt string `json:"prompt" validate:"required,min=5,max=500"`

	// AI model version
	ModelVersion samplerenums.ModelVersion `json:"modelVersion" validate:"required"`

	// Duration in seconds
	Duration int `json:"duration" validate:"required,min=1,max=30"`

	// Final audio format
	OutputFormat samplerenums.OutputFormat `json:"outputFormat" validate:"required"`

	SampleName string `json:"sampleName" validate:"required,min=1,max=500"`

	Email string `json:"email" validate:"required,email"`
}