package usecase

import "appointment-service/internal/model"

type PatientClient interface {
	GetAuthorizedPatient(authorizationHeader string) (*model.Patient, error)
}
