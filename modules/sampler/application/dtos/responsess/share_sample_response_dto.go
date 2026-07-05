package samplerresponsessdtos
type ShareSampleResponseDTO struct {
	Id              string  `json:"id"`
	SampleId        *string `json:"sampleId,omitempty"`        // omitempty por si viene nil
	SampleVersionId *string `json:"sampleVersionId,omitempty"` // omitempty por si viene nil
	UserId          string  `json:"userId"`
	Likes           int     `json:"likes"`
	Downloads       int     `json:"downloads"`
	CreatedAt       string  `json:"createdAt"`
}