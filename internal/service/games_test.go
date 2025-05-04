package service

import (
	"connectfour/internal/db"
	"connectfour/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func mockedGamesService() (*GamesService, *db.MockUserRepository, *db.MockGameRepository) {
	ur := db.MockUserRepository{}
	u := NewUserService(&ur, 0)
	sr := db.MockGameRepository{}
	return NewGamesService(u, &sr), &ur, &sr
}

func mockedGames() []model.Game {
	output := make([]model.Game, 3)
	output[0] = model.NewGame(user1, true)
	output[1] = model.NewGame(user2, true)
	// include one non-public game
	output[2] = model.NewGame(user1, false)
	return output
}

func TestGamesService_AllOpenGames(t *testing.T) {
	// Arrange
	list := mockedGames()
	s, ur, sr := mockedGamesService()
	ur.On("FindByEmail", mock.AnythingOfType("string")).Return(user1, nil)
	sr.On("List", mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(list, nil)

	// Act
	games := s.AllOpenGames(user1.Email)

	// Assert
	assert.Len(t, games, len(list)-1, "Expected number of games to match all Public games in the test set")
	assert.Equal(t, model.Created, games[0].Status)
	assert.Equal(t, model.Created, games[1].Status)
}

func TestGamesService_JoinGame_SavesToDb(t *testing.T) {
	list := mockedGames()
	game := &list[0]
	s, ur, sr := mockedGamesService()
	sr.On("Fetch", mock.AnythingOfType("string")).Return(*game, nil)
	sr.On("Save", mock.AnythingOfType("model.Game")).Return(true)
	ur.On("FindByEmail", mock.AnythingOfType("string")).Return(user2, nil)

	// Act
	err := s.JoinGame(game.Key, user2.Email)

	// Assert
	assert.NoError(t, err, "Expected no error when joining game")
}

func TestGamesService_AllMyGames(t *testing.T) {
	// Arrange
	list := mockedGames()
	s, ur, sr := mockedGamesService()
	ur.On("FindByEmail", mock.AnythingOfType("string")).Return(user1, nil)
	sr.On("List", user1.Id, mock.AnythingOfType("string")).Return(append([]model.Game{}, list[0], list[2]), nil)

	// Act
	games := s.AllMyGames(user1.Email)

	// Assert
	assert.Len(t, games, 2)
	sr.AssertCalled(t, "List", user1.Id, "")
}

func TestGamesService_CreateGame_SavesToDb(t *testing.T) {
	// Arrange
	s, ur, sr := mockedGamesService()
	sr.On("Save", mock.AnythingOfType("model.Game")).Return(true)
	ur.Mock.On("FindByEmail", mock.AnythingOfType("string")).Return(user1, nil)

	// Act
	resp := s.NewGame(user1.Email, true)

	// Assert
	sr.AssertCalled(t, "Save", mock.AnythingOfType("model.Game"))
	assert.Equal(t, model.Created, resp.Status)
	assert.Equal(t, user1.Email, resp.CreatedBy)
}
