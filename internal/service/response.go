package service

import (
	"connectfour/internal/model"
	"time"
)

// ErrorResponse represents an error returned by the API
// swagger:model
type ErrorResponse struct {
	// Error message
	// example: Game not found
	Error string `json:"error"`
}

// NewGameResponse represents a newly created game
// swagger:model
type NewGameResponse struct {
	// Unique identifier for the game
	// example: a1b2c3d4
	Key string `json:"key"`

	// Game creation time
	// example: 2023-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Email of the player who created the game
	// example: player1@example.com
	CreatedBy string `json:"created_by"`

	// Current status of the game
	// example: CREATED
	Status model.GameStatus `json:"status"`
}

func NewGameResponseFromGame(game model.Game) NewGameResponse {
	return NewGameResponse{
		Key:       game.Key,
		CreatedAt: game.CreatedAt,
		CreatedBy: game.Player1.Email,
		Status:    game.Status,
	}
}

// GameStateResponse represents the current state of a game
// swagger:model
//
//	@example: {
//		   "board": {
//						"1": "       ",
//					 	"2": "       ",
//					 	"3": "X  O   ",
//					 	"4": "X OX  O",
//				     	"5": "XOXOOOX",
//				     	"6": "XOXOXOX"
//		   },
//		   "key": "three-word-key",
//		   "player1_email": "player1@example.com",
//		   "player1_name": "Player 1",
//		   "player2_email": "player2@example.com",
//		   "player2_name": "Player 2",
//		   "player_turn": 1,
//		   "player_turn_email": "player1@example.com",
//		   "player_turn_name": "Player 1",
//		   "status": "started"
//		 }
type GameStateResponse struct {
	// Unique identifier for the game
	// example: three-word-key
	Key string `json:"key"`

	// Current status of the game
	// example: started
	// allowed:
	Status model.GameStatus `json:"status"`

	// Player whose turn it is (1 or 2)
	// example: 1
	PlayerTurn int `json:"player_turn"`

	// Name of the player whose turn it is
	// example: Player1
	PlayerTurnName string `json:"player_turn_name"`

	// Email of the player whose turn it is
	// example: player1@example.com
	PlayerTurnEmail string `json:"player_turn_email"`

	// Game board representation
	// example: {
	//			 	"1": "       ",
	//			 	"2": "       ",
	//			 	"3": "X  O   ",
	//			 	"4": "X OX  O",
	//		     	"5": "XOXOOOX",
	//		     	"6": "XOXOXOX"
	//		    }
	Board map[int]string `json:"board"`

	// Name of player 1
	// example: Player1
	Player1Name string `json:"player1_name"`

	// Name of player 2
	// example: Player2
	Player2Name string `json:"player2_name"`

	// Email of player 1
	// example: player1@example.com
	Player1Email string `json:"player1_email"`

	// Email of player 2
	// example: player2@example.com
	Player2Email string `json:"player2_email"`
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

// CreateUserResponse represents a newly created user
// swagger:model
type CreateUserResponse struct {
	// User's display name
	// example: Player1
	Name string `json:"name"`

	// User's email address
	// example: player1@example.com
	Email string `json:"email"`
}

func NewCreateUserResponse(u model.User) CreateUserResponse {
	return CreateUserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}
