package repository

import (
	"doctor-service/internal/model"
	"errors"
)

type InMemoryDoctorRepository struct {
	doctors map[string]*model.Doctor
	emails  map[string]string
}

func NewInMemoryDoctorRepository() *InMemoryDoctorRepository {
	return &InMemoryDoctorRepository{
		doctors: make(map[string]*model.Doctor),
		emails:  make(map[string]string),
	}
}

func (r *InMemoryDoctorRepository) Create(doctor *model.Doctor) error {
	if _, exists := r.emails[doctor.Email]; exists {
		return errors.New("doctor with this email already exists")
	}
	r.doctors[doctor.ID] = doctor
	r.emails[doctor.Email] = doctor.ID
	return nil
}

func (r *InMemoryDoctorRepository) GetByID(id string) (*model.Doctor, error) {
	doctor, exists := r.doctors[id]
	if !exists {
		return nil, errors.New("doctor not found")
	}
	return doctor, nil
}

func (r *InMemoryDoctorRepository) GetAll() ([]*model.Doctor, error) {
	var doctors []*model.Doctor
	for _, doctor := range r.doctors {
		doctors = append(doctors, doctor)
	}
	return doctors, nil
}
