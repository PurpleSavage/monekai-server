package commonmiddlewares

import "net/http"

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

