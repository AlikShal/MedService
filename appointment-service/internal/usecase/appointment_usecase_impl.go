package usecase

import (
	"appointment-service/internal/model"
	"appointment-service/internal/repository"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AppointmentUsecaseImpl struct {
	repo          repository.AppointmentRepository
	doctorClient  DoctorClient
	patientClient PatientClient
}

func NewAppointmentUsecaseImpl(repo repository.AppointmentRepository, doctorClient DoctorClient, patientClient PatientClient) *AppointmentUsecaseImpl {
	return &AppointmentUsecaseImpl{
		repo:          repo,
		doctorClient:  doctorClient,
		patientClient: patientClient,
	}
}

func (u *AppointmentUsecaseImpl) CreateAppointment(title, description, doctorID string, scheduledAt time.Time, authorizationHeader string) (*model.Appointment, error) {
	title = strings.TrimSpace(title)
	doctorID = strings.TrimSpace(doctorID)

	if title == "" {
		return nil, errors.New("title is required")
	}
	if doctorID == "" {
		return nil, errors.New("doctor_id is required")
	}
	if scheduledAt.IsZero() {
		return nil, errors.New("scheduled_at is required")
	}
	if scheduledAt.Before(time.Now().Add(-1 * time.Minute)) {
		return nil, errors.New("scheduled_at must be in the future")
	}

	doctor, err := u.doctorClient.CheckDoctorExists(doctorID)
	if err != nil {
		log.Printf("Failed to check doctor existence for doctor ID %s: %v", doctorID, err)
		return nil, errors.New("doctor does not exist or service unavailable")
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		log.Printf("Failed to resolve patient from token: %v", err)
		return nil, errors.New("patient profile not found or service unavailable")
	}

	now := time.Now().UTC()
	appointment := &model.Appointment{
		ID:          uuid.New().String(),
		Title:       title,
		Description: strings.TrimSpace(description),
		DoctorID:    doctorID,
		DoctorName:  doctor.FullName,
		PatientID:   patient.ID,
		PatientName: patient.FullName,
		ScheduledAt: scheduledAt.UTC(),
		Status:      model.StatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = u.repo.Create(appointment)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (u *AppointmentUsecaseImpl) GetAppointmentByID(id, role, authorizationHeader string) (*model.Appointment, error) {
	appointment, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(role, "admin") {
		return appointment, nil
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		return nil, errors.New("patient profile not found or service unavailable")
	}
	if appointment.PatientID != patient.ID {
		return nil, errors.New("appointment not found")
	}

	return appointment, nil
}

func (u *AppointmentUsecaseImpl) GetAllAppointments(role, authorizationHeader string) ([]*model.Appointment, error) {
	if strings.EqualFold(role, "admin") {
		return u.repo.GetAll()
	}

	patient, err := u.patientClient.GetAuthorizedPatient(authorizationHeader)
	if err != nil {
		return nil, errors.New("patient profile not found or service unavailable")
	}

	return u.repo.GetByPatientID(patient.ID)
}

func (u *AppointmentUsecaseImpl) UpdateAppointmentStatus(id string, status model.Status) error {
	current, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	switch status {
	case model.StatusNew, model.StatusInProgress, model.StatusDone, model.StatusCancelled:
	default:
		return errors.New("invalid status")
	}

	if current.Status == model.StatusDone || current.Status == model.StatusCancelled {
		return errors.New("cannot change a finalized appointment")
	}
	if current.Status == model.StatusNew && status == model.StatusDone {
		return errors.New("appointment must move to in_progress before done")
	}

	return u.repo.UpdateStatus(id, status)
}
