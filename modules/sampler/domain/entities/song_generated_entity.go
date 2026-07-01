package samplerentities

import (
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
)
type SongGeneratorResponse struct {
	GenerationId     string
	StatusGeneration samplerenums.Status
	UserId	uuid.UUID
}