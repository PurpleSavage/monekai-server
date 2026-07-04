package samplerresponsessdtos

type SampleResponseDTO struct {
	Id           string `json:"id"`
	SampleName   string `json:"sampleName"`
	Prompt       string `json:"prompt"`
	AudioUrl     *string `json:"audioUrl"`
	Duration     int    `json:"duration"`
	OutputFormat string `json:"outputFormat"`
	ModelVersion string `json:"modelVersion"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}