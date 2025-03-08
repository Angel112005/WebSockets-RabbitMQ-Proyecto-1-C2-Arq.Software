package persistence

import (
	"fmt"
	"log"
	"websocket-server/config"
	"github.com/rabbitmq/amqp091-go"
	"websocket-server/infrastructure/websocket"

	// "github.com/rabbitmq/amqp091-go"
)

// Suscribirse a RabbitMQ y enviar mensajes a WebSocket
func StartRabbitMQConsumer(hub *websocket.Hub) {
	msgs, err := config.RabbitMQChannel.Consume(
		"notificaciones.nuevas",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ No se pudo suscribir a la cola: %v", err)
	}

	go func() {
		for msg := range msgs {
			fmt.Println("📩 Mensaje recibido de RabbitMQ:", string(msg.Body))
			hub.Broadcast(msg.Body)
		}
	}()
}

// Escuchar mensajes de RabbitMQ y enviarlos por WebSocket
func ListenForNotifications(hub *websocket.Hub) {
	conn, err := amqp091.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ No se pudo abrir un canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notificaciones.nuevas", // Nombre de la cola
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("❌ No se pudo declarar la cola: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ No se pudo consumir mensajes: %v", err)
	}

	fmt.Println("📡 Escuchando eventos de notificaciones.nuevas...")

	// Escuchar mensajes y transmitirlos por WebSocket
	for msg := range msgs {
		fmt.Printf("📩 Mensaje recibido: %s\n", msg.Body)
		hub.Broadcast(msg.Body)
	}
}