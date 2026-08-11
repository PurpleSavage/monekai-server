package samplerports

type SamplerVersionPort interface{
	SaveEditSample()(error)
}