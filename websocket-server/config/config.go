
package config

import (
	"fmt"
	"log"
	// "os"

	"github.com/rabbitmq/amqp091-go"
)

// Configuración de RabbitMQ
var RabbitMQConn *amqp091.Connection
var RabbitMQChannel *amqp091.Channel

func InitRabbitMQ() {
	var err error
	RabbitMQConn, err = amqp091.Dial("amqp://agmc:112005@54.157.231.159:5672/")
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a RabbitMQ: %v", err)
	}

	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("❌ No se pudo abrir un canal en RabbitMQ: %v", err)
	}

	fmt.Println("✅ Conectado a RabbitMQ")
}

