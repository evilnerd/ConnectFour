package service

import (
	"connectfour/internal/db"
	"connectfour/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var (
	user1 = model.User{Name: "Dick", Email: "dick@evilnerd.nl", Token: "hunter2"}
	user2 = model.User{Name: "Sanae", Email: "sanae@evilnerd.nl", Token: "secret123"}
)

func init() {
	user1.Id = 1
	user2.Id = 2
}

func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	expected := user1

	repo := db.NewMockUserRepository()
	repo.On("Create", expected).Return(expected, nil)
	s := NewUserService(repo, time.Minute*5)

	// Act
	u1, err1 := s.CreateUser(expected.Email, expected.Name, expected.Token)
	u2, err2 := s.CreateUser("not.an.email", expected.Name, expected.Token)

	// Assert
	assert.EqualValues(t, expected, u1, "Expected the generated user to match the input values.")
	assert.NoError(t, err1, "Did not expect an error.")

	assert.Empty(t, u2, "Expected the generated user to be empty.")
	assert.Error(t, err2, "Expected an error, since the e-mail address was invalid")
}

func TestUserService_FindUserByEmail(t *testing.T) {

	// Arrange
	repo := db.NewMockUserRepository()
	repo.On("FindByEmail", "dick@evilnerd.nl").Return(user1, nil)
	repo.On("FindByEmail", mock.Anything).Return(model.User{}, errors.New("user not found"))

	s := NewUserService(repo, time.Minute*5)

	// Act
	u1, err1 := s.FindUserByEmail(user1.Email)
	u2, err2 := s.FindUserByEmail("other@evilnerd.nl")

	// Assert
	assert.EqualValues(t, u1, user1, "Expected the returned user to have the same values as the mock database user")
	assert.NoError(t, err1, "Expected the error to be nil")

	assert.Empty(t, u2, "Expected the returned user for the unknown address to be empty")
	assert.Error(t, err2, "Expected a User not found error.")
}

func TestUserService_Cache(t *testing.T) {
	// Arrange
	repo := db.NewMockUserRepository()
	repo.On("Create", user1).Panic("this should not be called")
	s := NewUserService(repo, time.Minute*5)

	// Act
	s.Cache(&user1)
	u1, err1 := s.FindUserByEmail(user1.Email)

	// Assert
	assert.NoError(t, err1, "Expected the error to be nil")
	assert.EqualValues(t, user1, u1, "Expected the returned user to match the input values")
	repo.AssertNotCalled(t, "FindByEmail")
}
