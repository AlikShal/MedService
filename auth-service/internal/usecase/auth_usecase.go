package usecase

import "auth-service/internal/model"

type AuthResult struct {
	Token string         `json:"token"`
	User  model.SafeUser `json:"user"`
}

type AuthUsecase interface {
	Register(fullName, email, password string) (*AuthResult, error)
	Login(email, password string) (*AuthResult, error)
	GetUserByID(id string) (*model.User, error)
}
