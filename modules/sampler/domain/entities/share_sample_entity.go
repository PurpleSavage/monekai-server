package samplerentities

import "github.com/google/uuid"

type ShareSampleEntity struct {
	Id              uuid.UUID
	SampleId        *uuid.UUID
	SampleVersionId *string
	UserId          string
	Likes           int
	Downloads       int
	CreatedAt       string
}