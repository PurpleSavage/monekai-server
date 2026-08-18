package samplerequestsdto

type UpdateEditedSampleURLDTO struct {
	FinalAudioURL string `json:"finalAudioUrl" validate:"required,url"`
}
