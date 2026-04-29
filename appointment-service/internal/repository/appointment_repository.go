package repository

import "appointment-service/internal/model"

type AppointmentRepository interface {
	Create(appointment *model.Appointment) error
	GetByID(id string) (*model.Appointment, error)
	GetAll() ([]*model.Appointment, error)
	GetByPatientID(patientID string) ([]*model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
}
