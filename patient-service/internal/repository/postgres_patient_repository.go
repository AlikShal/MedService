package repository

import (
	"database/sql"
	"errors"
	"patient-service/internal/model"
)

type PostgresPatientRepository struct {
	db *sql.DB
}

func NewPostgresPatientRepository(db *sql.DB) *PostgresPatientRepository {
	return &PostgresPatientRepository{db: db}
}

func (r *PostgresPatientRepository) Save(patient *model.Patient) (*model.Patient, error) {
	existing, err := r.GetByUserID(patient.UserID)
	if err != nil && err.Error() != "patient not found" {
		return nil, err
	}

	if existing == nil {
		query := `
			INSERT INTO patient_profiles (id, user_id, full_name, email, phone, date_of_birth, notes)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::date, $7)
			RETURNING created_at, updated_at
		`
		if err := r.db.QueryRow(
			query,
			patient.ID,
			patient.UserID,
			patient.FullName,
			patient.Email,
			patient.Phone,
			patient.DateOfBirth,
			patient.Notes,
		).Scan(&patient.CreatedAt, &patient.UpdatedAt); err != nil {
			return nil, err
		}
		return patient, nil
	}

	query := `
		UPDATE patient_profiles
		SET full_name = $2,
			email = $3,
			phone = $4,
			date_of_birth = NULLIF($5, '')::date,
			notes = $6,
			updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, created_at, updated_at
	`
	if err := r.db.QueryRow(
		query,
		patient.UserID,
		patient.FullName,
		patient.Email,
		patient.Phone,
		patient.DateOfBirth,
		patient.Notes,
	).Scan(&patient.ID, &patient.CreatedAt, &patient.UpdatedAt); err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PostgresPatientRepository) GetByUserID(userID string) (*model.Patient, error) {
	query := `
		SELECT id, user_id, full_name, email, phone,
		       COALESCE(TO_CHAR(date_of_birth, 'YYYY-MM-DD'), ''), notes, created_at, updated_at
		FROM patient_profiles
		WHERE user_id = $1
	`

	return r.scanPatient(r.db.QueryRow(query, userID))
}

func (r *PostgresPatientRepository) GetByID(id string) (*model.Patient, error) {
	query := `
		SELECT id, user_id, full_name, email, phone,
		       COALESCE(TO_CHAR(date_of_birth, 'YYYY-MM-DD'), ''), notes, created_at, updated_at
		FROM patient_profiles
		WHERE id = $1
	`

	return r.scanPatient(r.db.QueryRow(query, id))
}

func (r *PostgresPatientRepository) GetAll() ([]*model.Patient, error) {
	query := `
		SELECT id, user_id, full_name, email, phone,
		       COALESCE(TO_CHAR(date_of_birth, 'YYYY-MM-DD'), ''), notes, created_at, updated_at
		FROM patient_profiles
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := make([]*model.Patient, 0)
	for rows.Next() {
		patient := &model.Patient{}
		if err := rows.Scan(
			&patient.ID,
			&patient.UserID,
			&patient.FullName,
			&patient.Email,
			&patient.Phone,
			&patient.DateOfBirth,
			&patient.Notes,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, rows.Err()
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *PostgresPatientRepository) scanPatient(scanner rowScanner) (*model.Patient, error) {
	patient := &model.Patient{}
	err := scanner.Scan(
		&patient.ID,
		&patient.UserID,
		&patient.FullName,
		&patient.Email,
		&patient.Phone,
		&patient.DateOfBirth,
		&patient.Notes,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("patient not found")
	}
	if err != nil {
		return nil, err
	}
	return patient, nil
}
