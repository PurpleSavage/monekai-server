package samplerresponsessdtos

type EffectsResponseDTO struct {
	Reverb     *int  `json:"reverb,omitempty"`
	SlowPitch  *int  `json:"slowPitch,omitempty"`
	Saturation *int  `json:"saturation,omitempty"`
	Delay      *int  `json:"delay,omitempty"`
	LowPass    *int  `json:"lowPass,omitempty"`
	HighPass   *int  `json:"highPass,omitempty"`
	Gain       *int  `json:"gain,omitempty"`
	Reverse    *bool `json:"reverse,omitempty"`
}
