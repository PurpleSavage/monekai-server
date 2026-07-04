package commonresponsesdtos

type PaginatedResponse[T any] struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Data  []T`json:"data"`
}