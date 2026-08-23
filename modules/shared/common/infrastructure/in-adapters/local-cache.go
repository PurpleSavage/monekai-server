package commoninadapters

import (
	"strings"
	"sync"
	"time"

	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)
type DataCache[T any] struct {
	ttlSec int64
	data T
}

type LocalCache[T any] struct {
	mu      sync.RWMutex  
	Storage map[string]*DataCache[T]
}
func NewLocalCache[T any]() commonports.LocalCachePort[T] {
	return &LocalCache[T]{
		Storage: make(map[string]*DataCache[T]),
	}
}
func (l *LocalCache[T]) Set(data T, key string, ttlSecs int) (*T,error) {

	keyParsed := strings.TrimSpace(key)
	if keyParsed== "" {
		return nil , commondomainerrors.NewValidationError(
			"data void",
			"Null data cannot be saved.",
		)	
	}
	end := time.Now().Unix() + int64(ttlSecs)
	dataToSave := &DataCache[T]{
		ttlSec: end,
		data:   data,
	}
	
	l.mu.Lock()
	l.Storage[keyParsed] = dataToSave
	l.mu.Unlock()

	
	return &data, nil
	
}

func (l *LocalCache[T]) Get(key string) (*T, error) {
	keyParsed := strings.TrimSpace(key)
	if keyParsed== "" {
		return nil , commondomainerrors.NewValidationError(
			"data void",
			"Null data cannot be saved.",
		)	
	}
	l.mu.RLock()
	cachedData, ok := l.Storage[keyParsed]
	l.mu.RUnlock()
	if !ok || cachedData == nil {
		return nil, commondomainerrors.NewValidationError(
			"key",
			"Key not valid",
		)
	}
	timeNow := time.Now().Unix()
	if timeNow > cachedData.ttlSec {
		l.mu.Lock()
		delete(l.Storage, keyParsed)
		l.mu.Unlock()
		return nil, commondomainerrors.NewValidationError(
			"key_expired", 
			"The requested key has expired.",
		)
	}
	return &cachedData.data, nil
}