package repository

import (
	"appointment-service/internal/model"
	"database/sql"
	"errors"
)

type PostgresAppointmentRepository struct {
	db *sql.DB
}

func NewPostgresAppointmentRepository(db *sql.DB) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{db: db}
}

func (r *PostgresAppointmentRepository) Create(appointment *model.Appointment) error {
	query := `
		INSERT INTO appointments (
			id, title, description, doctor_id, doctor_name, patient_id, patient_name, scheduled_at, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		appointment.ID,
		appointment.Title,
		appointment.Description,
		appointment.DoctorID,
		appointment.DoctorName,
		appointment.PatientID,
		appointment.PatientName,
		appointment.ScheduledAt,
		appointment.Status,
	).Scan(&appointment.CreatedAt, &appointment.UpdatedAt)
}

func (r *PostgresAppointmentRepository) GetByID(id string) (*model.Appointment, error) {
	query := `
		SELECT id, title, description, doctor_id, doctor_name, patient_id, patient_name,
		       scheduled_at, status, created_at, updated_at
		FROM appointments
		WHERE id = $1
	`

	appointment := &model.Appointment{}
	err := r.db.QueryRow(query, id).Scan(
		&appointment.ID,
		&appointment.Title,
		&appointment.Description,
		&appointment.DoctorID,
		&appointment.DoctorName,
		&appointment.PatientID,
		&appointment.PatientName,
		&appointment.ScheduledAt,
		&appointment.Status,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("appointment not found")
	}
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *PostgresAppointmentRepository) GetAll() ([]*model.Appointment, error) {
	return r.list(`
		SELECT id, title, description, doctor_id, doctor_name, patient_id, patient_name,
		       scheduled_at, status, created_at, updated_at
		FROM appointments
		ORDER BY scheduled_at ASC, created_at DESC
	`)
}

func (r *PostgresAppointmentRepository) GetByPatientID(patientID string) ([]*model.Appointment, error) {
	return r.list(`
		SELECT id, title, description, doctor_id, doctor_name, patient_id, patient_name,
		       scheduled_at, status, created_at, updated_at
		FROM appointments
		WHERE patient_id = $1
		ORDER BY scheduled_at ASC, created_at DESC
	`, patientID)
}

func (r *PostgresAppointmentRepository) UpdateStatus(id string, status model.Status) error {
	result, err := r.db.Exec(`
		UPDATE appointments
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, id, status)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

func (r *PostgresAppointmentRepository) list(query string, args ...interface{}) ([]*model.Appointment, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := make([]*model.Appointment, 0)
	for rows.Next() {
		appointment := &model.Appointment{}
		if err := rows.Scan(
			&appointment.ID,
			&appointment.Title,
			&appointment.Description,
			&appointment.DoctorID,
			&appointment.DoctorName,
			&appointment.PatientID,
			&appointment.PatientName,
			&appointment.ScheduledAt,
			&appointment.Status,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}
