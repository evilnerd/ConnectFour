package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	player1 = NewUser("Sanae", "sanae@evilnerd.nl")
	player2 = NewUser("Dick", "dick@evilnerd.nl")
)

func TestGame_Join_OkIfStateIsCreated(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)

	// Act
	err := game.Join(player2)

	// Assert
	assert.Nil(t, err, "Expected the error to be empty")
}

func TestGame_Join_NotOkIfStateBeyondCreated(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)
	game.Status = Started
	// Act
	err := game.Join(player2)
	// Assert
	assert.Error(t, err, "Expected to get an error when joining a model that has already started")
}

func TestGame_Join_SetsStatusToStarted_AndSecondPlayerName(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)

	// Act
	err := game.Join(player2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, Started, game.Status, "Expected that the model status has been progressed to 'StartedAt'")
	assert.Equal(t, player2, game.Player2, "Expected the Player 2 name to be set")
	assert.Equal(t, 1, game.PlayerTurn)
	assert.NotNilf(t, game.StartedAt, "Expected the started time to be set")
}

func TestGame_Play_NotOkIfStateIsNotStarted(t *testing.T) {
	// Arrange
	game1 := NewGame(player1, true)
	game2 := NewGame(player1, true)
	game2.Status = Finished

	// Act
	err1 := game1.Play(player1, 1)
	err2 := game2.Play(player1, 1)

	// Assert
	assert.Error(t, err1, "Expected an error for model 1 since the model hasn't started yet")
	assert.Error(t, err2, "Expected an error for model 2 since the model has already finished")

}

func TestGame_Play_OkIfGameIsStarted(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)
	_ = game.Join(player2)

	// Act
	err := game.Play(player1, 1)

	// Assert
	assert.Nil(t, err, "Expected no error since the model should be in the right state to play.")
}

func TestGame_Play_SetsOtherPlayerMove(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)
	_ = game.Join(player2)

	// Act
	_ = game.Play(player1, 1)

	// Assert
	assert.Equal(t, 2, game.PlayerTurn)
}
