package service

import (
	"connectfour/internal/db"
	"connectfour/internal/model"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/mail"
	"strings"
	"time"
)

type UserService struct {
	repo      db.UserRepository
	userCache *Cache[string, *model.User]
}

type UserExistsError struct {
	email string
}

func (e UserExistsError) Error() string {
	return fmt.Sprintf("user with email %s already exists", e.email)
}

func NewUserExistsError(email string) UserExistsError {
	return UserExistsError{email: email}
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

	log.Debugf("Creating user %s (%s)...", name, email)
	email = strings.ToLower(email)
	if !validateEmail(email) {
		return model.User{}, errors.New("invalid email address")
	}

	user, err := s.FindUserByEmail(email)
	if err != nil {
		log.Errorf("Could not determine if user %s already exists: %v", email, err)
		return model.User{}, fmt.Errorf("could not determine if the user already exists")
	}
	if !user.Empty() {
		return model.User{}, NewUserExistsError(email)
	}

	user = model.User{
		Name:  name,
		Email: strings.ToLower(email),
		Token: token,
	}
	user, err = s.repo.Create(user)
	if err != nil {
		log.Errorf("Error creating user %s (%s): %v", name, email, err)
		return model.User{}, err
	}
	log.Debugf("Succeeded creating user %s (%s)...", name, email)
	return user, nil
}

// validateEmail returns true when the e-mail address is valid.
func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
