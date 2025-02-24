package service

import (
	"connectfour/internal/model"
	"time"
)

type NewGameResponse struct {
	Key       string           `json:"key"`
	CreatedAt time.Time        `json:"created_at"`
	CreatedBy string           `json:"created_by"`
	Status    model.GameStatus `json:"status"`
}

func NewGameResponseFromGame(game model.Game) NewGameResponse {
	return NewGameResponse{
		Key:       game.Key,
		CreatedAt: game.CreatedAt,
		CreatedBy: game.Player1.Email,
		Status:    game.Status,
	}
}

type GameStateResponse struct {
	Key             string           `json:"key"`
	Status          model.GameStatus `json:"status"`
	PlayerTurn      int              `json:"player_turn"` // either 1 or 2
	PlayerTurnName  string           `json:"player_turn_name"`
	PlayerTurnEmail string           `json:"player_turn_email"`
	Board           map[int]string   `json:"board"`
	Player1Name     string           `json:"player1_name"`
	Player2Name     string           `json:"player2_name"`
	Player1Email    string           `json:"player1_email"`
	Player2Email    string           `json:"player2_email"`
}

func NewGameStateResponse(game model.Game) GameStateResponse {
	return GameStateResponse{
		Key:             game.Key,
		Status:          game.Status,
		PlayerTurn:      game.PlayerTurn,
		PlayerTurnEmail: game.CurrentPlayer().Email,
		PlayerTurnName:  game.CurrentPlayer().Name,
		Board:           game.Board.Map(),
		Player1Name:     game.Player1.Name,
		Player2Name:     game.Player2.Name,
		Player1Email:    game.Player1.Email,
		Player2Email:    game.Player2.Email,
	}
}

type CreateUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewCreateUserResponse(u model.User) CreateUserResponse {
	return CreateUserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}
