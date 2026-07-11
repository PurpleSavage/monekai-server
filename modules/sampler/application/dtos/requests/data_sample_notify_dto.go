package samplerequestsdto

import samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"


type DataSampleNotify struct {
	UserID   string
	SampleID string
	Data     *samplerresponsessdtos.SampleResponseDTO // nil en el caso de error
}