package usecase

import "chat-service/internal/model"

type AppointmentClient interface {
	GetAppointment(id string, authorizationHeader string) (*model.Appointment, error)
}
