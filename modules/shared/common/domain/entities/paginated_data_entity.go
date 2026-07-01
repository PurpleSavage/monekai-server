package commonentities
type PaginatedResult[T any] struct {
	Total int
	Limit int
	Page  int
	Data  []T 
}