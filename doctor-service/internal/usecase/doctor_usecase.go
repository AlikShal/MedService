package usecase

import "doctor-service/internal/model"

type DoctorUsecase interface {
	CreateDoctor(fullName, specialization, email, office string) (*model.Doctor, error)
	GetDoctorByID(id string) (*model.Doctor, error)
	GetAllDoctors() ([]*model.Doctor, error)
}
