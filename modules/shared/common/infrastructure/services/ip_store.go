package commonservices
import (
	"sync"
	"time"
	"golang.org/x/time/rate"
)
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type IPStore struct {
	mu      sync.Mutex
	clients map[string]*client
}
func NewIPStore() *IPStore {
	store := &IPStore{clients: make(map[string]*client)}
	// Hilo secundario para limpiar IPs inactivas y evitar fugas de memoria
	go store.cleanupRoutine()
	return store
}

func (s *IPStore) cleanupRoutine() {
	for {
		time.Sleep(time.Minute)
		s.mu.Lock()
		for ip, c := range s.clients {
			// Si el cliente no ha hecho peticiones en 3 minutos, lo borramos de la RAM
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(s.clients, ip)
			}
		}
		s.mu.Unlock()
	}
}

func (s *IPStore) Allow(ip string) bool {
	s.mu.Lock()
	c, exists := s.clients[ip]
	if !exists {
		// 5 peticiones por segundo, con ráfaga de 10
		c = &client{limiter: rate.NewLimiter(rate.Limit(5), 10)}
		s.clients[ip] = c
	}
	c.lastSeen = time.Now()
	s.mu.Unlock()

	// Retorna si la petición es permitida o no
	return c.limiter.Allow()
}
