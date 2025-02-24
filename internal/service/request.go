package service

type NewGameRequest struct {
	Public bool `json:"public"`
}

type PlayMoveRequest struct {
	Column int `json:"column"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
