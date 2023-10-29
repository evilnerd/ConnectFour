package server

import (
	"connectfour/game"
	"errors"
	"time"
)

type GameStatus string

type Game struct {
	Key         string
	Player1Name string // red
	Player2Name string // yellow
	PlayerTurn  int    // either 1 or 2
	CreatedAt   time.Time
	StartedAt   time.Time
	FinishedAt  time.Time

	Public bool
	Status GameStatus
	board  *game.Board
}

const (
	Created  GameStatus = "created"
	Started  GameStatus = "started"
	Finished GameStatus = "finished"
	Aborted  GameStatus = "aborted"
)

// NewGame will create a new game, add it to the list and set the status to 'created'. It needs a second
// player to start.
func NewGame(player1Name string, public bool) Game {
	g := Game{
		Key:         GenerateKey(3),
		Player1Name: player1Name,
		Public:      public,
		CreatedAt:   time.Now(),
		Status:      Created,
		board:       &game.Board{},
	}
	games[g.Key] = &g
	return g
}

// Join will add the second player to the game and set the status to 'started'.
func (g *Game) Join(player2Name string) error {
	if g.Status != Created {
		return errors.New("you can only join a game that has status 'Created'")
	}

	g.Status = Started
	g.StartedAt = time.Now()
	g.Player2Name = player2Name
	g.PlayerTurn = 1

	return nil
}

// Play will make a play for the current player on the specified column, and set the other player's turn
// unless the game has ended.
// Column is 1-based (so acceptable values are 1-7)
func (g *Game) Play(column int) error {
	if g.Status == Created {
		return errors.New("this game is not started yet, still waiting for the second player")
	}

	if g.Status == Finished || g.Status == Aborted {
		return errors.New("this game has finished and you can't play any more moves on it")
	}

	if !g.board.AddDisc(column-1, g.playerDisc()) {
		return errors.New("invalid move")
	}

	if g.board.HasConnectFour() {
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

func (g *Game) playerDisc() game.Disc {
	if g.PlayerTurn == 1 {
		return game.RedDisc
	}
	return game.YellowDisc
}

// CurrentPlayerName returns the name of the player whose turn it is.
func (g *Game) CurrentPlayerName() string {
	if g.PlayerTurn == 1 {
		return g.Player1Name
	}
	return g.Player2Name
}
