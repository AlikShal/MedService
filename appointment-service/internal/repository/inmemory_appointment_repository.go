package repository

import (
	"appointment-service/internal/model"
	"errors"
	"time"
)

type InMemoryAppointmentRepository struct {
	appointments map[string]*model.Appointment
}

func NewInMemoryAppointmentRepository() *InMemoryAppointmentRepository {
	return &InMemoryAppointmentRepository{
		appointments: make(map[string]*model.Appointment),
	}
}

func (r *InMemoryAppointmentRepository) Create(appointment *model.Appointment) error {
	r.appointments[appointment.ID] = appointment
	return nil
}

func (r *InMemoryAppointmentRepository) GetByID(id string) (*model.Appointment, error) {
	appointment, exists := r.appointments[id]
	if !exists {
		return nil, errors.New("appointment not found")
	}
	return appointment, nil
}

func (r *InMemoryAppointmentRepository) GetAll() ([]*model.Appointment, error) {
	var appointments []*model.Appointment
	for _, appointment := range r.appointments {
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *InMemoryAppointmentRepository) GetByPatientID(patientID string) ([]*model.Appointment, error) {
	var appointments []*model.Appointment
	for _, appointment := range r.appointments {
		if appointment.PatientID == patientID {
			appointments = append(appointments, appointment)
		}
	}
	return appointments, nil
}

func (r *InMemoryAppointmentRepository) UpdateStatus(id string, status model.Status) error {
	appointment, exists := r.appointments[id]
	if !exists {
		return errors.New("appointment not found")
	}
	appointment.Status = status
	appointment.UpdatedAt = time.Now()
	return nil
}
