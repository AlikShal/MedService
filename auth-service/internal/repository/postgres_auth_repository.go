package repository

import (
	"auth-service/internal/model"
	"database/sql"
	"errors"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO auth_users (id, full_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, user.ID, user.FullName, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresAuthRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at
		FROM auth_users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresAuthRepository) GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at
		FROM auth_users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
