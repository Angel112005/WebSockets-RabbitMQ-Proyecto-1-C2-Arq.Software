package api

import (
	"encoding/json"
	"log"
	"my-notification-api/application"
	"my-notification-api/domain"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	createNotification *application.CreateNotification
}

func NewNotificationController(cn *application.CreateNotification) *NotificationController {
	return &NotificationController{createNotification: cn}
}

// Función auxiliar para intentar parsear la fecha en distintos formatos
func parseFecha(fechaStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04",          // Sin segundos
		"2006-01-02T15:04:05",       // Sin zona horaria
		"2006-01-02T15:04:05Z07:00", // RFC3339
	}

	var t time.Time
	var err error

	for _, format := range formats {
		t, err = time.Parse(format, fechaStr)
		if err == nil {
			return t, nil
		}
	}

	return t, err
}

func (c *NotificationController) Create(ctx *gin.Context) {
	// Definir estructura para recibir datos del JSON
	var input struct {
		PatientName     string `json:"patient_name"`
		DoctorID        int    `json:"doctor_id"`
		AppointmentDate string `json:"appointment_date"`
		Status          string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Intentar parsear la fecha con la función personalizada
	parsedDate, err := parseFecha(input.AppointmentDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}

	// Crear la notificación con la fecha parseada correctamente
	notification := domain.Notification{
		PatientName:     input.PatientName,
		DoctorID:        input.DoctorID,
		AppointmentDate: parsedDate,
		Status:          input.Status,
	}

	// Guardar en la base de datos
	err = c.createNotification.Execute(notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la notificación"})
		return
	}

	// 📢 Enviar notificación a RabbitMQ
	publicarNotificacionRabbitMQ(notification)

	ctx.JSON(http.StatusCreated, gin.H{"message": "Notificación guardada correctamente"})
}

// 📌 Función para enviar notificación a la cola "notificaciones.nuevas"
func publicarNotificacionRabbitMQ(notification domain.Notification) {
	conn, err := amqp.Dial("amqp://agmc:112005@54.157.231.159:5672/")
	if err != nil {
		log.Fatalf("❌ Error al conectar con RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error al abrir canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notificaciones.nuevas", // 🆕 Cola para notificaciones
		true,                     // Durable
		false,                    // Auto-delete
		false,                    // Exclusiva
		false,                    // No-wait
		nil,                      // Argumentos
	)
	if err != nil {
		log.Fatalf("❌ Error al declarar la cola: %s", err)
	}

	body, _ := json.Marshal(notification)
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing Key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("❌ Error al publicar mensaje en RabbitMQ: %s", err)
	}

	log.Printf("📢 Evento enviado a RabbitMQ (notificaciones.nuevas): %s", body)
}




// package api

// import (
// 	"my-notification-api/application"
// 	"my-notification-api/domain"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// type NotificationController struct {
// 	createNotification *application.CreateNotification
// }

// func NewNotificationController(cn *application.CreateNotification) *NotificationController {
// 	return &NotificationController{createNotification: cn}
// }

// // Función auxiliar para intentar parsear la fecha en distintos formatos
// func parseFecha(fechaStr string) (time.Time, error) {
// 	formats := []string{
// 		"2006-01-02T15:04",          // Sin segundos
// 		"2006-01-02T15:04:05",       // Sin zona horaria
// 		"2006-01-02T15:04:05Z07:00", // RFC3339
// 	}

// 	var t time.Time
// 	var err error

// 	for _, format := range formats {
// 		t, err = time.Parse(format, fechaStr)
// 		if err == nil {
// 			return t, nil
// 		}
// 	}

// 	return t, err
// }

// func (c *NotificationController) Create(ctx *gin.Context) {
// 	// Definir estructura para recibir datos del JSON
// 	var input struct {
// 		PatientName     string `json:"patient_name"`
// 		DoctorID        int    `json:"doctor_id"`
// 		AppointmentDate string `json:"appointment_date"`
// 		Status          string `json:"status"`
// 	}

// 	if err := ctx.ShouldBindJSON(&input); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
// 		return
// 	}

// 	// Intentar parsear la fecha con la función personalizada
// 	parsedDate, err := parseFecha(input.AppointmentDate)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
// 		return
// 	}

// 	// Crear la notificación con la fecha parseada correctamente
// 	notification := domain.Notification{
// 		PatientName:     input.PatientName,
// 		DoctorID:        input.DoctorID,
// 		AppointmentDate: parsedDate,
// 		Status:          input.Status,
// 	}

// 	// Guardar en la base de datos
// 	err = c.createNotification.Execute(notification)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la notificación"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{"message": "Notificación guardada correctamente"})
// }
