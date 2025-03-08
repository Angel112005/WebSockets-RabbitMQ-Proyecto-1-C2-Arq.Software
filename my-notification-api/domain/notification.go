package domain

import "time"

type Notification struct {
	ID             int       `json:"id"`
	PatientName    string    `json:"patient_name"`
	DoctorID       int       `json:"doctor_id"`
	AppointmentDate time.Time `json:"appointment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
