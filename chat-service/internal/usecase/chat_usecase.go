package usecase

import "chat-service/internal/model"

type ChatUsecase interface {
	CreateThread(appointmentID, subject, authorizationHeader string, claims *Claims) (*model.ChatThread, error)
	ListThreads(authorizationHeader string, claims *Claims) ([]*model.ChatThread, error)
	GetMessages(threadID, authorizationHeader string, claims *Claims) ([]*model.ChatMessage, error)
	SendMessage(threadID, body, authorizationHeader string, claims *Claims) (*model.ChatMessage, error)
}
