package communityresponsesdtos

import "github.com/google/uuid"

type SharedByDTO struct {
	ID    uuid.UUID `json:"userId"` // Mapeado explícitamente a userId para el Front
	Email string    `json:"email"`
}