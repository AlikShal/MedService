package usecase

import (
	"appointment-service/internal/model"
	"time"
)

type AppointmentUsecase interface {
	CreateAppointment(title, description, doctorID string, scheduledAt time.Time, authorizationHeader string) (*model.Appointment, error)
	GetAppointmentByID(id, role, authorizationHeader string) (*model.Appointment, error)
	GetAllAppointments(role, authorizationHeader string) ([]*model.Appointment, error)
	UpdateAppointmentStatus(id string, status model.Status) error
}
