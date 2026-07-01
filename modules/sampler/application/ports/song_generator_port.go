package samplerports

import (
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
)
type SongGeneratorPort interface {
	GenerateSong(vo *samplervalueobjects.PayloadGenerateSongVO) (*samplerentities.SongGeneratorResponse, error)
} 