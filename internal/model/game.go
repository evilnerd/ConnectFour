package model

import (
	"errors"
	"strings"
	"time"
)

type GameStatus string

type Game struct {
	Key        string
	Player1    User
	Player2    User
	PlayerTurn int // either 1 or 2
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time

	Public bool
	Status GameStatus
	Board  Board
}

const (
	Created  GameStatus = "created"
	Started  GameStatus = "started"
	Finished GameStatus = "finished"
	Aborted  GameStatus = "aborted"
	Unknown  GameStatus = "unknown" // used in the client to indicate the status should be fetched.
)

// NewGame will create a new game, add it to the list and set the status to 'created'. It needs a second
// player to start.
func NewGame(player1 User, public bool) Game {
	g := Game{
		Key:       GenerateKey(3),
		Player1:   player1,
		Public:    public,
		CreatedAt: time.Now(),
		Status:    Created,
		Board:     Board{},
	}
	return g
}

// Join will add the second player to the game and set the status to 'started'.
func (g *Game) Join(joining User) error {
	if g.Status != Created && joining.Id != g.Player2.Id && joining.Id != g.Player1.Id {
		return errors.New("you can only join a game that has status 'Created'")
	}

	// If the joining player is already part of the game, then
	// just do nothing to the game's state.
	if g.Player1.Is(joining) || g.Player2.Is(joining) {
		return nil
	}

	// if the second player wasn't set yet, then now is the time
	// to start the game.
	if g.Player2.Empty() {
		g.Player2 = joining
		g.PlayerTurn = 1
		g.Status = Started
		g.StartedAt = time.Now()
	}

	return nil
}

// Play will make a play for the current player on the specified column, and set the other player's turn
// unless the game has ended.
// Column is 1-based (so acceptable values are 1-7)
func (g *Game) Play(user User, column int) error {

	if g.Status == Created {
		return errors.New("this game is not started yet, still waiting for the second player")
	}

	if g.Status == Finished || g.Status == Aborted {
		return errors.New("this game has finished and you can't play any more moves on it")
	}

	if !g.IsPlayerTurn(user.Email) {
		return errors.New("it is not your turn")
	}

	if !g.Board.AddDisc(column-1, g.playerDisc()) {
		return errors.New("invalid move")
	}

	if g.Board.HasConnectFour() {
		g.Status = Finished
		g.FinishedAt = time.Now()
	} else {
		g.switchPlayer()
	}

	return nil
}

func (g *Game) switchPlayer() {
	if g.PlayerTurn == 1 {
		g.PlayerTurn = 2
	} else {
		g.PlayerTurn = 1
	}
}

func (g *Game) playerDisc() Disc {
	if g.PlayerTurn == 1 {
		return RedDisc
	}
	return YellowDisc
}

// CurrentPlayer returns a pointer to the current player 'User'.
func (g *Game) CurrentPlayer() *User {
	if g.PlayerTurn == 1 {
		return &g.Player1
	}
	return &g.Player2
}

// CurrentPlayerEmail returns the email of the player whose turn it is.
func (g *Game) CurrentPlayerEmail() string {
	return g.CurrentPlayer().Email
}

func (g *Game) IsPlayerTurn(email string) bool {
	return strings.EqualFold(g.CurrentPlayerEmail(), email)
}
