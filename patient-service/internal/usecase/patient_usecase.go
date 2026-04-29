package usecase

import "patient-service/internal/model"

type PatientUsecase interface {
	SaveProfile(userID, email, fullName, phone, dateOfBirth, notes string) (*model.Patient, error)
	GetByUserID(userID string) (*model.Patient, error)
	GetByID(id string) (*model.Patient, error)
	GetAll() ([]*model.Patient, error)
}
