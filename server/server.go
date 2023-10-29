package server

var (
	games map[string]*Game
)

func init() {
	games = make(map[string]*Game)
}

func GetGame(key string) *Game {
	return games[key]
}

func GetGameState(key string) GameStateResponse {
	return NewGameStateResponse(*games[key])
}

func AllOpenGames() []NewGameResponse {
	responses := make([]NewGameResponse, 0)
	for _, game := range games {
		if game.Public {
			responses = append(responses, NewGameResponseFromGame(*game))
		}
	}

	return responses
}

func GameExists(key string) bool {
	_, ok := games[key]
	return ok
}
