package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame_Join_OkIfStateIsCreated(t *testing.T) {
	// Arrange
	game := NewGame("Sanae", true)

	// Act
	err := game.Join("Dick")

	// Assert
	assert.Nil(t, err, "Expected the error to be empty")
}

func TestGame_Join_NotOkIfStateBeyondCreated(t *testing.T) {
	// Arrange
	game := NewGame("Sanae", true)
	game.Status = Started
	// Act
	err := game.Join("Dick")
	// Assert
	assert.Error(t, err, "Expected to get an error when joining a game that has already started")
}

func TestGame_Join_SetsStatusToStarted_AndSecondPlayerName(t *testing.T) {
	// Arrange
	game := NewGame("Sanae", true)

	// Act
	game.Join("Dick")

	// Assert
	assert.Equal(t, Started, game.Status, "Expected that the game status has been progressed to 'StartedAt'")
	assert.Equal(t, "Dick", game.Player2Name, "Expected the Player 2 name to be set")
	assert.Equal(t, 1, game.PlayerTurn)
	assert.NotNilf(t, game.StartedAt, "Expected the started time to be set")
}

func TestGame_Play_NotOkIfStateIsNotStarted(t *testing.T) {
	// Arrange
	game1 := NewGame("Sanae", true)
	game2 := NewGame("Also Sanae", true)
	game2.Status = Finished

	// Act
	err1 := game1.Play(1)
	err2 := game2.Play(1)

	// Assert
	assert.Error(t, err1, "Expected an error for game 1 since the game hasn't started yet")
	assert.Error(t, err2, "Expected an error for game 2 since the game has already finished")

}

func TestGame_Play_OkIfGameIsStarted(t *testing.T) {
	// Arrange
	game := NewGame("Sanae", true)
	game.Join("Dick")

	// Act
	err := game.Play(1)

	// Assert
	assert.Nil(t, err, "Expected no error since the game should be in the right state to play.")
}

func TestGame_Play_SetsOtherPlayerMove(t *testing.T) {
	// Arrange
	game := NewGame("Sanae", true)
	game.Join("Dick")

	// Act
	game.Play(1)

	// Assert
	assert.Equal(t, 2, game.PlayerTurn)
}
