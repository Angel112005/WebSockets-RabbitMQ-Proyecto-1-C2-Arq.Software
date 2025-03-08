package api

import (
	"fmt"
	gorilla "github.com/gorilla/websocket"
	"net/http"
	ws "websocket-server/infrastructure/websocket"
)

// Upgrader para WebSocket
var upgrader = gorilla.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(hub *ws.Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil) // Gorilla upgrade
        if err != nil {
            fmt.Println("❌ Error al conectar WebSocket:", err)
            return
        }

        // Usar tu Hub
        hub.AddClient(conn)
        fmt.Println("🟢 Cliente WebSocket conectado")
    }
}