package commonutils

import "github.com/google/uuid"
func UUIDPtrToStringPtr(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	s := id.String()
	return &s
}