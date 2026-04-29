package usecase

import (
	"doctor-service/internal/model"
	"doctor-service/internal/repository"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type DoctorUsecaseImpl struct {
	repo repository.DoctorRepository
}

func NewDoctorUsecaseImpl(repo repository.DoctorRepository) *DoctorUsecaseImpl {
	return &DoctorUsecaseImpl{repo: repo}
}

func (u *DoctorUsecaseImpl) CreateDoctor(fullName, specialization, email, office string) (*model.Doctor, error) {
	fullName = strings.TrimSpace(fullName)
	email = strings.ToLower(strings.TrimSpace(email))

	if fullName == "" {
		return nil, errors.New("full_name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	doctor := &model.Doctor{
		ID:             uuid.New().String(),
		FullName:       fullName,
		Specialization: strings.TrimSpace(specialization),
		Email:          email,
		Office:         strings.TrimSpace(office),
	}
	err := u.repo.Create(doctor)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (u *DoctorUsecaseImpl) GetDoctorByID(id string) (*model.Doctor, error) {
	return u.repo.GetByID(id)
}

func (u *DoctorUsecaseImpl) GetAllDoctors() ([]*model.Doctor, error) {
	return u.repo.GetAll()
}
