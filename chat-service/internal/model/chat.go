package model

import "time"

type ThreadStatus string

const (
	ThreadStatusOpen   ThreadStatus = "open"
	ThreadStatusClosed ThreadStatus = "closed"
)

type ChatThread struct {
	ID            string       `json:"id"`
	AppointmentID string       `json:"appointment_id"`
	PatientID     string       `json:"patient_id"`
	PatientName   string       `json:"patient_name"`
	Subject       string       `json:"subject"`
	Status        ThreadStatus `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	LastMessageAt *time.Time   `json:"last_message_at,omitempty"`
}

type ChatMessage struct {
	ID           string    `json:"id"`
	ThreadID     string    `json:"thread_id"`
	SenderUserID string    `json:"sender_user_id"`
	SenderRole   string    `json:"sender_role"`
	SenderName   string    `json:"sender_name"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
}

type Patient struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Appointment struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PatientID   string `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Status      string `json:"status"`
}
