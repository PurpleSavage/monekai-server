package samplerentities

import "github.com/google/uuid"

type ShareSampleEntity struct {
	Id              uuid.UUID
	SampleId        *uuid.UUID
	SampleVersionId *uuid.UUID
	UserId          string
	Likes           int
	Downloads       int
	CreatedAt       string
}