package model

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

type Appointment struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DoctorID    string    `json:"doctor_id"`
	DoctorName  string    `json:"doctor_name"`
	PatientID   string    `json:"patient_id"`
	PatientName string    `json:"patient_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
