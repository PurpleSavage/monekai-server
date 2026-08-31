package samplerequestsdto

type PresignedURLRequestDTO struct {
	Key  string `json:"key" validate:"required,oneof=originals editions"`
	Name string `json:"name" validate:"required,min=1,max=255"`
}
