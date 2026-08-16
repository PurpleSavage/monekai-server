package samplerequestsdto

type EffectsDTO struct {
	Reverb     *int  `json:"reverb"`
	SlowPitch  *int  `json:"slowPitch"`
	Saturation *int  `json:"saturation"`
	Delay      *int  `json:"delay"`
	LowPass    *int  `json:"lowPass"`
	HighPass   *int  `json:"highPass"`
	Gain       *int  `json:"gain"`
	Reverse    *bool `json:"reverse"`
}