package notificationsevents

import (
	"sync"
)


type SSEEvent struct {
	Name string `json:"event"`
	Data string `json:"data"`
}

type SSEManager struct {
	mu      sync.RWMutex
	clients map[string]chan SSEEvent
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]chan SSEEvent),
	}
}

// Register: El usuario entra o recarga la página
func (m *SSEManager) Register(email string) chan SSEEvent {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Si recargó la página, sobreescribimos la referencia en el mapa.
	// Al NO cerrar el canal viejo, si un webhook venía viajando en ese milisegundo,
	// escribirá en un canal huérfano que luego el Garbage Collector destruirá.
	ch := make(chan SSEEvent, 10)
	m.clients[email] = ch
	return ch
}

// Unregister: El usuario se desconectó (cerró pestaña o recargó)
func (m *SSEManager) Unregister(email string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//  sacamos del mapa para liberar la clave.
	delete(m.clients, email)
}

// BroadcastToUser: Envía el evento si el cliente está conectado
func (m *SSEManager) BroadcastToUser(email string, eventName string, data string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Buscamos si el usuario tiene un canal activo en el mapa
	if ch, exists := m.clients[email]; exists {
		select {
		case ch <- SSEEvent{Name: eventName, Data: data}:
			// Éxito: El evento se encoló en el buffer
		default:
			// Si el cliente tiene el buffer lleno (pestaña congelada),
			// pasamos de largo para no bloquear el hilo del Webhook.
		}
	}
}