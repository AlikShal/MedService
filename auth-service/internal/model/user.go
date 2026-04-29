package model

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RolePatient Role = "patient"
)

type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type SafeUser struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Sanitize() SafeUser {
	return SafeUser{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
