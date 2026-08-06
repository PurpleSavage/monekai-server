package commonmiddlewares

import (
	"net"
	"net/http"
	commonservices "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/services"
)



type ThrotlerRequest struct{}

func NewThrotlerRequest()*ThrotlerRequest{
	return &ThrotlerRequest{}
}
func (t *ThrotlerRequest) RateLimit(store *commonservices.IPStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Simplemente llamamos al método público del servicio
		if !store.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}