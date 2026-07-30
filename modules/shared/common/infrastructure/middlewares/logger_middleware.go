package commonmiddlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (l *LoggerMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.NewString()[:8]
		start := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		elapsed := time.Since(start).Round(time.Millisecond)
		methodColored := colorForMethod(r.Method)
		statusColored := colorForStatus(rec.Status)
		cyan := "\033[36m"
		reset := "\033[0m"

		log.Println(fmt.Sprintf(
			"%s[%s]%s %s%s%s %s%d%s %s%v%s | %s %s%s%s",
			cyan, reqID, reset,
			methodColored, r.Method, reset,
			statusColored, rec.Status, reset,
			cyan, elapsed, reset,
			r.RemoteAddr,
			cyan, r.URL.Path, reset,
		))
	})
}

func colorForStatus(status int) string {
	switch {
	case status >= 500:
		return "\033[31m"
	case status >= 400:
		return "\033[38;5;208m"
	case status >= 300:
		return "\033[36m"
	default:
		return "\033[32m"
	}
}

func colorForMethod(method string) string {
	switch method {
	case http.MethodGet:
		return "\033[36m"
	case http.MethodPost:
		return "\033[32m"
	case http.MethodPatch:
		return "\033[33m"
	case http.MethodDelete:
		return "\033[31m"
	default:
		return "\033[37m"
	}
}