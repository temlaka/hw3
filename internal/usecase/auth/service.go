package auth

import (
	"errors"
	"fmt"
	"hw3/internal/domain"
	"hw3/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   UserRepository
	secret string
}

func NewService(repo UserRepository, secret string) *Service {
	return &Service{
		repo:   repo,
		secret: secret,
	}
}

func (s *Service) Register(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", domain.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {

		return "", fmt.Errorf("password too long: %w", domain.ErrInvalidInput)

	} else if err != nil {
		return "", fmt.Errorf("hash password: %w", err)

	}
	newUser := domain.User{
		Username:     username,
		PasswordHash: hash,
	}

	t, err := token.GenerateToken(username, s.secret)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	if err := s.repo.Create(newUser); err != nil {
		return "", err
	}
	return t, nil
}

func (s *Service) Login(username string, password string) (string, error) {

	user, err := s.repo.GetByUsername(username)

	if errors.Is(err, domain.ErrUserNotFound) {
		return "", domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", domain.ErrInvalidCredentials
	} else if err != nil {
		return "", fmt.Errorf("compare password: %w", err)
	}

	t, err := token.GenerateToken(username, s.secret)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return t, nil
}
