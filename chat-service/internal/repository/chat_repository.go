package repository

import "chat-service/internal/model"

type ChatRepository interface {
	CreateThread(thread *model.ChatThread) error
	GetThreadByID(id string) (*model.ChatThread, error)
	GetThreadByAppointmentID(appointmentID string) (*model.ChatThread, error)
	ListThreads() ([]*model.ChatThread, error)
	ListThreadsByPatientID(patientID string) ([]*model.ChatThread, error)
	AddMessage(message *model.ChatMessage) error
	ListMessages(threadID string) ([]*model.ChatMessage, error)
}
