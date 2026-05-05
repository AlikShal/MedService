package usecase

import "chat-service/internal/model"

type PatientClient interface {
	GetAuthorizedPatient(authorizationHeader string) (*model.Patient, error)
}
