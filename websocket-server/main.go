package main

import (
	"fmt"
	"net/http"
	"websocket-server/config"
	"websocket-server/infrastructure/persistence"
	"websocket-server/infrastructure/websocket"
	"websocket-server/server"
)

func main() {
	config.InitRabbitMQ()

	hub := websocket.NewHub()
	server.SetupRoutes(hub)

	go persistence.StartRabbitMQConsumer(hub)

	fmt.Println("🚀 WebSocket server corriendo en :8082")
	http.ListenAndServe(":8082", nil)
}
