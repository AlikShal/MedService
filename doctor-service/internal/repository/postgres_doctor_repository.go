package repository

import (
	"database/sql"
	"doctor-service/internal/model"
	"errors"
)

type PostgresDoctorRepository struct {
	db *sql.DB
}

func NewPostgresDoctorRepository(db *sql.DB) *PostgresDoctorRepository {
	return &PostgresDoctorRepository{db: db}
}

func (r *PostgresDoctorRepository) Create(doctor *model.Doctor) error {
	query := `
		INSERT INTO doctors (id, full_name, specialization, email, office)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	return r.db.QueryRow(
		query,
		doctor.ID,
		doctor.FullName,
		doctor.Specialization,
		doctor.Email,
		doctor.Office,
	).Scan(&doctor.CreatedAt)
}

func (r *PostgresDoctorRepository) GetByID(id string) (*model.Doctor, error) {
	query := `
		SELECT id, full_name, specialization, email, office, created_at
		FROM doctors
		WHERE id = $1
	`

	doctor := &model.Doctor{}
	err := r.db.QueryRow(query, id).Scan(
		&doctor.ID,
		&doctor.FullName,
		&doctor.Specialization,
		&doctor.Email,
		&doctor.Office,
		&doctor.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("doctor not found")
	}
	if err != nil {
		return nil, err
	}

	return doctor, nil
}

func (r *PostgresDoctorRepository) GetAll() ([]*model.Doctor, error) {
	query := `
		SELECT id, full_name, specialization, email, office, created_at
		FROM doctors
		ORDER BY full_name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]*model.Doctor, 0)
	for rows.Next() {
		doctor := &model.Doctor{}
		if err := rows.Scan(
			&doctor.ID,
			&doctor.FullName,
			&doctor.Specialization,
			&doctor.Email,
			&doctor.Office,
			&doctor.CreatedAt,
		); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	return doctors, rows.Err()
}
