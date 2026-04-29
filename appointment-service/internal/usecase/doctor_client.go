package usecase

import "appointment-service/internal/model"

type DoctorClient interface {
	CheckDoctorExists(doctorID string) (*model.Doctor, error)
}
