package persistence

import (
	"database/sql"
	"my-notification-api/domain"
)

type NotificationMySQL struct {
	db *sql.DB
}

func NewNotificationMySQL(db *sql.DB) *NotificationMySQL {
	return &NotificationMySQL{db: db}
}

func (r *NotificationMySQL) Save(notification domain.Notification) error {
	_, err := r.db.Exec(
		"INSERT INTO notifications (patient_name, doctor_id, appointment_date, status) VALUES (?, ?, ?, ?)",
		notification.PatientName, notification.DoctorID, notification.AppointmentDate, notification.Status,
	)
	return err
}
