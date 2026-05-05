package repository

import (
	"chat-service/internal/model"
	"database/sql"
	"errors"
)

type PostgresChatRepository struct {
	db *sql.DB
}

func NewPostgresChatRepository(db *sql.DB) *PostgresChatRepository {
	return &PostgresChatRepository{db: db}
}

func (r *PostgresChatRepository) CreateThread(thread *model.ChatThread) error {
	query := `
		INSERT INTO chat_threads (id, appointment_id, patient_id, patient_name, subject, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at, last_message_at
	`

	return r.db.QueryRow(
		query,
		thread.ID,
		thread.AppointmentID,
		thread.PatientID,
		thread.PatientName,
		thread.Subject,
		thread.Status,
	).Scan(&thread.CreatedAt, &thread.UpdatedAt, &thread.LastMessageAt)
}

func (r *PostgresChatRepository) GetThreadByID(id string) (*model.ChatThread, error) {
	return r.scanThread(r.db.QueryRow(threadSelectQuery()+" WHERE id = $1", id))
}

func (r *PostgresChatRepository) GetThreadByAppointmentID(appointmentID string) (*model.ChatThread, error) {
	return r.scanThread(r.db.QueryRow(threadSelectQuery()+" WHERE appointment_id = $1", appointmentID))
}

func (r *PostgresChatRepository) ListThreads() ([]*model.ChatThread, error) {
	return r.listThreads(threadSelectQuery() + " ORDER BY COALESCE(last_message_at, created_at) DESC")
}

func (r *PostgresChatRepository) ListThreadsByPatientID(patientID string) ([]*model.ChatThread, error) {
	return r.listThreads(
		threadSelectQuery()+" WHERE patient_id = $1 ORDER BY COALESCE(last_message_at, created_at) DESC",
		patientID,
	)
}

func (r *PostgresChatRepository) AddMessage(message *model.ChatMessage) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO chat_messages (id, thread_id, sender_user_id, sender_role, sender_name, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`
	if err := tx.QueryRow(
		query,
		message.ID,
		message.ThreadID,
		message.SenderUserID,
		message.SenderRole,
		message.SenderName,
		message.Body,
	).Scan(&message.CreatedAt); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE chat_threads
		SET updated_at = NOW(), last_message_at = $2
		WHERE id = $1
	`, message.ThreadID, message.CreatedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresChatRepository) ListMessages(threadID string) ([]*model.ChatMessage, error) {
	rows, err := r.db.Query(`
		SELECT id, thread_id, sender_user_id, sender_role, sender_name, body, created_at
		FROM chat_messages
		WHERE thread_id = $1
		ORDER BY created_at ASC
	`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*model.ChatMessage, 0)
	for rows.Next() {
		message := &model.ChatMessage{}
		if err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.SenderUserID,
			&message.SenderRole,
			&message.SenderName,
			&message.Body,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, rows.Err()
}

func threadSelectQuery() string {
	return `
		SELECT id, appointment_id, patient_id, patient_name, subject, status, created_at, updated_at, last_message_at
		FROM chat_threads
	`
}

func (r *PostgresChatRepository) listThreads(query string, args ...interface{}) ([]*model.ChatThread, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*model.ChatThread, 0)
	for rows.Next() {
		thread := &model.ChatThread{}
		if err := rows.Scan(
			&thread.ID,
			&thread.AppointmentID,
			&thread.PatientID,
			&thread.PatientName,
			&thread.Subject,
			&thread.Status,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.LastMessageAt,
		); err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, rows.Err()
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *PostgresChatRepository) scanThread(scanner rowScanner) (*model.ChatThread, error) {
	thread := &model.ChatThread{}
	err := scanner.Scan(
		&thread.ID,
		&thread.AppointmentID,
		&thread.PatientID,
		&thread.PatientName,
		&thread.Subject,
		&thread.Status,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.LastMessageAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("chat thread not found")
	}
	if err != nil {
		return nil, err
	}
	return thread, nil
}
