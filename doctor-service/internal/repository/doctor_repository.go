package repository

import "doctor-service/internal/model"

type DoctorRepository interface {
	Create(doctor *model.Doctor) error
	GetByID(id string) (*model.Doctor, error)
	GetAll() ([]*model.Doctor, error)
}
