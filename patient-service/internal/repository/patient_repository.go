package repository

import "patient-service/internal/model"

type PatientRepository interface {
	Save(patient *model.Patient) (*model.Patient, error)
	GetByUserID(userID string) (*model.Patient, error)
	GetByID(id string) (*model.Patient, error)
	GetAll() ([]*model.Patient, error)
}
