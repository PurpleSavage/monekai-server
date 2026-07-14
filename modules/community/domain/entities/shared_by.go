package communityentities

import "github.com/google/uuid"
type SharedBy struct {
	ID        uuid.UUID
	Name      string
	Email     string
}