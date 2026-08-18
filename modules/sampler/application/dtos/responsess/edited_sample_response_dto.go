package samplerresponsessdtos

type EditedSampleResponseDTO struct {
	ID            string             `json:"id"`
	Sample        SampleResponseDTO  `json:"sample"`
	Effects       EffectsResponseDTO `json:"effects"`
	FinalAudioURL string             `json:"finalAudioUrl"`
	CreatedAt     string             `json:"createdAt"`
}
