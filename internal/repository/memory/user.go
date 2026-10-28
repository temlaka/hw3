package memory

import (
	"hw3/internal/domain"
	"sync"
)

type UserRepo struct {
	storage map[string]domain.User
	mu      sync.RWMutex
}

func NewUserRepo() *UserRepo {

	return &UserRepo{
		storage: make(map[string]domain.User),
	}
}

func (u *UserRepo) Create(user domain.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.storage[user.Username]; ok {
		return domain.ErrUserExists
	}
	u.storage[user.Username] = user
	return nil

}

func (u *UserRepo) GetByUsername(username string) (domain.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	user, ok := u.storage[username]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil

}
