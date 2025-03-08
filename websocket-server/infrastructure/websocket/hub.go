package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// Hub maneja todas las conexiones activas
type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// Nueva instancia del Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

// Añadir un cliente
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

// Eliminar un cliente
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

// Enviar notificación a todos los clientes
func (h *Hub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("❌ Error enviando mensaje:", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
