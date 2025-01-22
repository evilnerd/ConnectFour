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
	output := make([]model.Game, 0)
	g1 := model.NewGame(user1, true)
	output = append(output, g1)
	g2 := model.NewGame(user2, true)
	output = append(output, g2)
	// include one non-public game
	g3 := model.NewGame(user1, false)
	output = append(output, g3)
	return output
}

func TestGamesService_AllOpenGames(t *testing.T) {
	// Arrange
	list := mockedGames()
	s, ur, sr := mockedGamesService()
	ur.Mock.On("FindByEmail", mock.AnythingOfType("string")).Return(user1, nil)
	sr.Mock.On("List", mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(list, nil)

	// Act
	games := s.AllOpenGames(user1.Email)

	// Assert
	assert.Len(t, games, len(list)-1, "Expected number of games to match all Public games in the test set")
	assert.Equal(t, model.Created, games[0].Status)
	assert.Equal(t, model.Created, games[1].Status)
}

func TestGamesService_JoinGame_SavesToDb(t *testing.T) {
	// s, ur, sr := mockedGamesService()
}
