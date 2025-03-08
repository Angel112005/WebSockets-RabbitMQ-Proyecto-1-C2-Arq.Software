package server

import (
	"net/http"
	"websocket-server/infrastructure/api"
	"websocket-server/infrastructure/websocket"
)

// Configurar rutas
func SetupRoutes(hub *websocket.Hub) {
	http.HandleFunc("/ws", api.WebSocketHandler(hub))
}
