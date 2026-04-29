package model

import "time"

type Doctor struct {
	ID             string    `json:"id"`
	FullName       string    `json:"full_name"`
	Specialization string    `json:"specialization"`
	Email          string    `json:"email"`
	Office         string    `json:"office"`
	CreatedAt      time.Time `json:"created_at"`
}
