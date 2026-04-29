package usecase

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseImpl struct {
	repo      repository.AuthRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthUsecase(repo repository.AuthRepository, jwtSecret string, tokenTTL time.Duration) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (u *AuthUsecaseImpl) Register(fullName, email, password string) (*AuthResult, error) {
	fullName = strings.TrimSpace(fullName)
	email = strings.ToLower(strings.TrimSpace(email))

	if fullName == "" {
		return nil, errors.New("full_name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	if _, err := u.repo.GetUserByEmail(email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.NewString(),
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         model.RolePatient,
	}
	if err := u.repo.CreateUser(user); err != nil {
		return nil, err
	}

	createdUser, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return u.buildResult(createdUser)
}

func (u *AuthUsecaseImpl) Login(email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return u.buildResult(user)
}

func (u *AuthUsecaseImpl) GetUserByID(id string) (*model.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *AuthUsecaseImpl) buildResult(user *model.User) (*AuthResult, error) {
	token, err := GenerateToken(user, u.jwtSecret, u.tokenTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user.Sanitize(),
	}, nil
}
