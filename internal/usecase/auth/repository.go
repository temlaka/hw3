package auth

import "hw3/internal/domain"

type UserRepository interface {
	Create(domain.User) error
	GetByUsername(username string) (domain.User, error)
}
