package usecase

import (
	"errors"
	"patient-service/internal/model"
	"patient-service/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type PatientUsecaseImpl struct {
	repo repository.PatientRepository
}

func NewPatientUsecase(repo repository.PatientRepository) *PatientUsecaseImpl {
	return &PatientUsecaseImpl{repo: repo}
}

func (u *PatientUsecaseImpl) SaveProfile(userID, email, fullName, phone, dateOfBirth, notes string) (*model.Patient, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user_id is required")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email is required")
	}
	if strings.TrimSpace(fullName) == "" {
		return nil, errors.New("full_name is required")
	}

	patient := &model.Patient{
		ID:          uuid.NewString(),
		UserID:      userID,
		FullName:    strings.TrimSpace(fullName),
		Email:       strings.ToLower(strings.TrimSpace(email)),
		Phone:       strings.TrimSpace(phone),
		DateOfBirth: strings.TrimSpace(dateOfBirth),
		Notes:       strings.TrimSpace(notes),
	}

	return u.repo.Save(patient)
}

func (u *PatientUsecaseImpl) GetByUserID(userID string) (*model.Patient, error) {
	return u.repo.GetByUserID(userID)
}

func (u *PatientUsecaseImpl) GetByID(id string) (*model.Patient, error) {
	return u.repo.GetByID(id)
}

func (u *PatientUsecaseImpl) GetAll() ([]*model.Patient, error) {
	return u.repo.GetAll()
}
