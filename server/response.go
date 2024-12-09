package server

import "time"

type NewGameResponse struct {
	Key       string     `json:"key"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	Status    GameStatus `json:"status"`
}

func NewGameResponseFromGame(game Game) NewGameResponse {
	return NewGameResponse{
		Key:       game.Key,
		CreatedAt: game.CreatedAt,
		CreatedBy: game.Player1Name,
		Status:    game.Status,
	}
}

type GameStateResponse struct {
	Key            string         `json:"key"`
	Status         GameStatus     `json:"status"`
	PlayerTurn     int            `json:"player_turn"` // either 1 or 2
	PlayerTurnName string         `json:"player_turn_name"`
	Board          map[int]string `json:"board"`
	Player1Name    string         `json:"player1Name"`
	Player2Name    string         `json:"player2Name"`
}

func NewGameStateResponse(game Game) GameStateResponse {
	return GameStateResponse{
		Key:            game.Key,
		Status:         game.Status,
		PlayerTurn:     game.PlayerTurn,
		PlayerTurnName: game.CurrentPlayerName(),
		Board:          game.board.Map(),
		Player1Name:    game.Player1Name,
		Player2Name:    game.Player2Name,
	}
}
