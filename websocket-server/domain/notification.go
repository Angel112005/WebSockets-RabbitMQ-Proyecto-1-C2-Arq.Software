package domain

type Notification struct {
	ID             int    `json:"id"`
	PatientName    string `json:"patient_name"`
	DoctorID       int    `json:"doctor_id"`
	AppointmentDate string `json:"appointment_date"`
	Status         string `json:"status"`
}
