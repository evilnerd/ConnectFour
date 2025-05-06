package service

// NewGameRequest represents the request to create a new game
// swagger:model
type NewGameRequest struct {
	// Whether the game should be public or private
	// example: true
	Public bool `json:"public"`
}

// PlayMoveRequest represents the request to make a move in a game
// swagger:model
type PlayMoveRequest struct {
	// Column where the player wants to place their disc (0-indexed)
	// Required: true
	// minimum: 0
	// maximum: 6
	// example: 3
	Column int `json:"column"`
}

// LoginRequest represents the user login credentials
// swagger:model
type LoginRequest struct {
	// User's email address
	// Required: true
	// example: player@example.com
	Email string `json:"email"`

	// User's password
	// Required: true
	// example: password123
	Password string `json:"password"`
}

// RegisterRequest represents the user registration data
// swagger:model
type RegisterRequest struct {
	// User's display name
	// Required: true
	// example: Player1
	Name string `json:"name"`

	// User's email address
	// Required: true
	// example: player@example.com
	Email string `json:"email"`

	// User's password
	// Required: true
	// example: password123
	Password string `json:"password"`
}
