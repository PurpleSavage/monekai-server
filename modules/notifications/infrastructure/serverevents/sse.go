package notificationsevents

import (
	"fmt"
	"net/http"
	"time"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/go-chi/chi/v5"
)

type AudioSSEHandler struct {
	sseManager     *SSEManager	
	authMiddleware *authmiddlewares.AuthMiddleware
}

// 2. El constructor independiente
func NewAudioSSEHandler(
	sm *SSEManager,
	am *authmiddlewares.AuthMiddleware,
) *AudioSSEHandler {
	return &AudioSSEHandler{
		sseManager:     sm,
		authMiddleware: am,
	}
}

// 3. El método que maneja el streaming de eventos
func (h *AudioSSEHandler) StreamSongStatus(w http.ResponseWriter, r *http.Request) {
	// Extraemos el email del contexto de forma segura gracias a tu middleware de JWT
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		http.Error(w, "Unauthorized: Session data not found in context", http.StatusUnauthorized)
		return
	}
	dto, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		http.Error(w, "Internal Server Error: Could not parse session data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	userChan := h.sseManager.Register(dto.Email)
	defer h.sseManager.Unregister(dto.Email)

	keepAliveTicker := time.NewTicker(15 * time.Second)
	defer keepAliveTicker.Stop()

	clientGone := r.Context().Done()

	for {
		select {
		case <-clientGone:
			return

		case <-keepAliveTicker.C:
			fmt.Fprintf(w, ": keep-alive\n\n")
			flusher.Flush()

		case eventData := <-userChan:
			if eventData.Name != "" {
				fmt.Fprintf(w, "event: %s\n", eventData.Name)
			}
			fmt.Fprintf(w, "data: %s\n\n", eventData)
			flusher.Flush()
		}
	}
}

// 4. Su propia función para mapear las rutas de este Handler específico
func MapSSERoutes(r chi.Router, h *AudioSSEHandler) {
	
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.RefreshToken) 
		r.Get("/sse/stream", h.StreamSongStatus)
	})
	
}