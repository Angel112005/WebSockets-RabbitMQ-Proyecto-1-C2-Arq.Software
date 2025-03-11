package main

import (
	"bytes"
	"encoding/json"

	// "fmt"
	"io/ioutil"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 📌 Estructura del mensaje recibido de RabbitMQ
type Cita struct {
	ID              int    `json:"id"`
	PatientName     string `json:"patient_name"`
	DoctorID        int    `json:"doctor_id"`
	AppointmentDate string `json:"appointment_date"`
	Status          string `json:"status"`
}

// 📌 Función para conectar a RabbitMQ y consumir mensajes
func consumirCitas() {
	conn, err := amqp.Dial("amqp://agmc:112005@54.157.231.159:5672/")
	if err != nil {
		log.Fatalf("❌ Error conectando a RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error abriendo canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"citas.nuevas", // Nombre de la cola
		true,           // Durable
		false,          // AutoDelete
		false,          // Exclusiva
		false,          // NoWait
		nil,            // Argumentos
	)
	if err != nil {
		log.Fatalf("❌ Error declarando la cola: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // Nombre de la cola
		"",     // Consumer
		true,   // AutoAck
		false,  // Exclusivo
		false,  // No-local
		false,  // NoWait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("❌ Error registrando el consumidor: %s", err)
	}

	// Procesar los mensajes de la cola
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("📥 Mensaje recibido de RabbitMQ: %s", d.Body)

			// Convertir el mensaje JSON a una estructura Cita
			var cita Cita
			err := json.Unmarshal(d.Body, &cita)
			if err != nil {
				log.Printf("❌ Error al parsear JSON: %s", err)
				continue
			}

			// Enviar la cita a API 2
			enviarCitaAPI2(cita)
		}
	}()

	log.Println("📡 Esperando mensajes de RabbitMQ...")
	<-forever
}

// 📌 Función para enviar la cita a API 2
func enviarCitaAPI2(cita Cita) {

	url := "http://localhost:8081/notificaciones" // API 2

	// Convertir la cita a JSON
	citaJSON, err := json.Marshal(cita)
	if err != nil {
		log.Printf("❌ Error al convertir la cita a JSON: %s", err)
		return
	}

	// Hacer la petición HTTP a API 2
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(citaJSON))
	if err != nil {
		log.Printf("❌ Error al enviar la cita a API 2: %s", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("📤 Respuesta de API 2: %s", string(body))
}

func main() {
	consumirCitas()
}
 
