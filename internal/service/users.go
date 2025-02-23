package service

import (
	"connectfour/internal/db"
	"connectfour/internal/model"
	"errors"
	"net/mail"
	"strings"
	"time"
)

type UserService struct {
	repo      db.UserRepository
	userCache *Cache[string, *model.User]
}

func NewUserService(repo db.UserRepository, cacheTtl time.Duration) *UserService {
	return &UserService{
		repo:      repo,
		userCache: NewCache[string, *model.User](cacheTtl),
	}
}

func (s UserService) Cache(user *model.User) {
	s.userCache.Store(strings.ToLower(user.Email), user)
}

func (s UserService) FindUserByEmail(email string) (model.User, error) {

	// Normalize
	email = strings.ToLower(email)

	// first look in the cache
	user, ok := s.userCache.Load(email)
	if !ok {
		user, err := s.repo.FindByEmail(email)
		if err == nil {
			s.userCache.Store(email, &user)
		}
		return user, err
	}
	return *user, nil
}

func (s UserService) CreateUser(email string, name string, token string) (model.User, error) {

	if !validateEmail(email) {
		return model.User{}, errors.New("invalid email address")
	}

	user := model.User{
		Name:  name,
		Email: email,
		Token: token,
	}
	user, err := s.repo.Create(user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// validateEmail returns true when the e-mail address is valid.
func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
