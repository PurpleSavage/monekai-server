package commonports

type LocalCachePort[T any] interface {
	Set(data T, key string, ttlSecs int) (*T, error)
	Get(key string) (*T, error)
}