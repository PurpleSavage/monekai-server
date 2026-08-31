package samplerresponsessdtos

type PresignedURLResponseDTO struct {
	UploadURL string `json:"uploadUrl"`
	Key       string `json:"key"`
}
