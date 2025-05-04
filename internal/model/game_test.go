package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	player1 = NewUser("Sanae", "sanae@evilnerd.nl")
	player2 = NewUser("Dick", "dick@evilnerd.nl")
	player3 = NewUser("Lucy", "lucy@evilnerd.nl")
)

// Give the test users a proper ID, since it's used in some test cases to compare users.
func init() {
	player1.Id = 1
	player2.Id = 2
	player3.Id = 3
}

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
	_ = game.Join(player2)
	game.Status = Started

	// Act
	err := game.Join(player3)
	// Assert
	assert.Error(t, err, "Expected to get an error when joining a game that has already started")
}

func TestGame_Join_SetsStatusToStarted_AndSecondPlayerName(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)

	// Act
	err := game.Join(player2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, Started, game.Status, "Expected that the game status has been progressed to 'StartedAt'")
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
	assert.Error(t, err1, "Expected an error for game 1 since the game hasn't started yet")
	assert.Error(t, err2, "Expected an error for game 2 since the game has already finished")

}

func TestGame_Play_OkIfGameIsStarted(t *testing.T) {
	// Arrange
	game := NewGame(player1, true)
	_ = game.Join(player2)

	// Act
	err := game.Play(player1, 1)

	// Assert
	assert.Nil(t, err, "Expected no error since the game should be in the right state to play.")
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
