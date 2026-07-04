package samplerequestsdto
type SampleGenerationInput struct {
	Prompt                    string  `json:"prompt"`
	ModelVersion              string  `json:"model_version"`
	Duration                  int     `json:"duration"`
	NormalizationStrategy     string  `json:"normalization_strategy"`
	ClassifierFreeGuidance    float32 `json:"classifier_free_guidance"`
	OutputFormat              string  `json:"output_format"`
	TopK                      int     `json:"top_k"`
	TopP                      float32 `json:"top_p"`
	Temperature               float32 `json:"temperature"`
}

