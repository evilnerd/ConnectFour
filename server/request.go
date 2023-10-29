package server

type NewGameRequest struct {
	Player1Name string `json:"player1_name"`
	Public      bool   `json:"public"`
}

type JoinGameRequest struct {
	Player2Name string `json:"player2_name"`
}

type PlayMoveRequest struct {
	Column int `json:"column"`
}
