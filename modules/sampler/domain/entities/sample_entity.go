package samplerentities

import (
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
)

type SampleEntity struct {
	Id           uuid.UUID
	SampleName   string
	Prompt       string
	AudioUrl     *string
	Duration     int
	OutputFormat samplerenums.OutputFormat
	ModelVersion samplerenums.ModelVersion
	Status samplerenums.Status
	CreatedAt string
}
