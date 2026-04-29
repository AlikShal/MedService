package model

import "time"

type Patient struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	DateOfBirth string    `json:"date_of_birth"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
