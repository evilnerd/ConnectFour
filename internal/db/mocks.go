package db

import (
	"connectfour/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) Create(u model.User) (model.User, error) {
	args := m.Called(u)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

type MockGameRepository struct {
	mock.Mock
}

func NewMockGameRepository() *MockGameRepository {
	return &MockGameRepository{}
}

func (m *MockGameRepository) Save(game model.Game) bool {
	args := m.Called(game)
	return args.Bool(0)
}

func (m *MockGameRepository) Fetch(key string) (model.Game, error) {
	args := m.Called(key)
	return args.Get(0).(model.Game), args.Error(1)
}

func (m *MockGameRepository) List(userId int64, status string) ([]model.Game, error) {
	args := m.Called(userId, status)
	return args.Get(0).([]model.Game), args.Error(1)
}
