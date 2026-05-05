package usecase

import (
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ChatUsecaseImpl struct {
	repo              repository.ChatRepository
	patientClient     PatientClient
	appointmentClient AppointmentClient
}

func NewChatUsecaseImpl(repo repository.ChatRepository, patientClient PatientClient, appointmentClient AppointmentClient) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		repo:              repo,
		patientClient:     patientClient,
		appointmentClient: appointmentClient,
	}
}

func (u *ChatUsecaseImpl) CreateThread(appointmentID, subject, authorizationHeader string, claims *Claims) (*model.ChatThread, error) {
	if claims == nil {
		return nil, errors.New("missing auth context")
	}
	if !strings.EqualFold(claims.Role, "patient") {
		return nil, errors.New("only patients can open appointment chat threads")
	}

	appointmentID = strings.TrimSpace(appointmentID)
	if appointmentID == "" {
		return nil, errors.New("appointment_id is required")
	}

	if existing, err := u.repo.GetThreadByAppointmentID(appointmentID); err == nil {
		if err := u.ensureThreadAccess(existing, authorizationHeader, claims); err != nil {
			return nil, err
		}
		return existing, nil
	} else if !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return nil, err
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		return nil, errors.New("patient profile not found or service unavailable")
	}

	appointment, err := u.appointmentClient.GetAppointment(appointmentID, authorizationHeader)
	if err != nil {
		return nil, err
	}
	if appointment.PatientID != patient.ID {
		return nil, errors.New("appointment not found")
	}

	subject = strings.TrimSpace(subject)
	if subject == "" {
		subject = appointment.Title
	}

	now := time.Now().UTC()
	thread := &model.ChatThread{
		ID:            uuid.New().String(),
		AppointmentID: appointment.ID,
		PatientID:     patient.ID,
		PatientName:   patient.FullName,
		Subject:       subject,
		Status:        model.ThreadStatusOpen,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := u.repo.CreateThread(thread); err != nil {
		return nil, err
	}

	return thread, nil
}

func (u *ChatUsecaseImpl) ListThreads(authorizationHeader string, claims *Claims) ([]*model.ChatThread, error) {
	if claims == nil {
		return nil, errors.New("missing auth context")
	}
	if strings.EqualFold(claims.Role, "admin") {
		return u.repo.ListThreads()
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		return nil, errors.New("patient profile not found or service unavailable")
	}

	return u.repo.ListThreadsByPatientID(patient.ID)
}

func (u *ChatUsecaseImpl) GetMessages(threadID, authorizationHeader string, claims *Claims) ([]*model.ChatMessage, error) {
	thread, err := u.repo.GetThreadByID(strings.TrimSpace(threadID))
	if err != nil {
		return nil, err
	}
	if err := u.ensureThreadAccess(thread, authorizationHeader, claims); err != nil {
		return nil, err
	}

	return u.repo.ListMessages(thread.ID)
}

func (u *ChatUsecaseImpl) SendMessage(threadID, body, authorizationHeader string, claims *Claims) (*model.ChatMessage, error) {
	thread, err := u.repo.GetThreadByID(strings.TrimSpace(threadID))
	if err != nil {
		return nil, err
	}
	if err := u.ensureThreadAccess(thread, authorizationHeader, claims); err != nil {
		return nil, err
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.New("message body is required")
	}
	if len(body) > 2000 {
		return nil, errors.New("message body must be 2000 characters or fewer")
	}

	senderName := claims.Email
	if strings.EqualFold(claims.Role, "patient") {
		patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
		if err != nil {
			return nil, errors.New("patient profile not found or service unavailable")
		}
		senderName = patient.FullName
	} else if strings.EqualFold(claims.Role, "admin") {
		senderName = "Clinic Administrator"
	}

	message := &model.ChatMessage{
		ID:           uuid.New().String(),
		ThreadID:     thread.ID,
		SenderUserID: claims.UserID,
		SenderRole:   strings.ToLower(claims.Role),
		SenderName:   senderName,
		Body:         body,
		CreatedAt:    time.Now().UTC(),
	}
	if err := u.repo.AddMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (u *ChatUsecaseImpl) ensureThreadAccess(thread *model.ChatThread, authorizationHeader string, claims *Claims) error {
	if claims == nil {
		return errors.New("missing auth context")
	}
	if strings.EqualFold(claims.Role, "admin") {
		return nil
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		return errors.New("patient profile not found or service unavailable")
	}
	if thread.PatientID != patient.ID {
		return errors.New("chat thread not found")
	}

	return nil
}
