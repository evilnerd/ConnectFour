package service

import (
	"connectfour/internal/db"
	"connectfour/internal/model"
	log "github.com/sirupsen/logrus"
)

type GamesService struct {
	userService    *UserService
	gameRepository db.GameRepository
}

func NewGamesService(userService *UserService, gamesRepository db.GameRepository) *GamesService {
	return &GamesService{
		userService:    userService,
		gameRepository: gamesRepository,
	}
}

func init() {
}

func (s GamesService) GetGame(key string) model.Game {
	game, err := s.gameRepository.Fetch(key)
	if err != nil {
		log.Errorf("Error fetching game from db: %v", err)
		return model.Game{}
	}
	return game
}

func (s GamesService) GetGameState(key string) GameStateResponse {
	game := s.GetGame(key)
	if game.Key != key {
		return GameStateResponse{}
	}
	return NewGameStateResponse(game)
}

func (s GamesService) AllOpenGames(email string) []NewGameResponse {
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		log.Errorf("Error finding user by email '%s': %v", email, err)
		return []NewGameResponse{}
	}

	games, err := s.gameRepository.List(user.Id, string(model.Created))
	if err != nil {
		log.Errorf("Error getting games: %v\n", err)
		return []NewGameResponse{}
	}

	output := make([]NewGameResponse, 0)
	for _, game := range games {
		if game.Public {
			output = append(output, NewGameResponseFromGame(game))
		}
	}
	return output
}

func (s GamesService) AllMyGames(email string) []NewGameResponse {
	log.Debug("Listing all games for user %s", email)
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		log.Errorf("Error finding user by email '%s': %v", email, err)
		return []NewGameResponse{}
	}

	if email == "" || user.Empty() {
		log.Errorf("No active user.")
		return []NewGameResponse{}
	}

	output := make([]NewGameResponse, 0)
	games, err := s.gameRepository.List(user.Id, "")
	if err != nil {
		log.Errorf("Error getting games: %v\n", err)
	}
	for _, game := range games {
		output = append(output, NewGameResponseFromGame(game))
	}
	return output
}

func (s GamesService) GameExists(key string) bool {
	game, err := s.gameRepository.Fetch(key)
	if err != nil {
		log.Errorf("Error determining if game '%s' exists or not: %v", key, err)
		return false
	}
	return game.Key == key
}

func (s GamesService) JoinGame(key string, player2Email string) error {
	user, err := s.userService.FindUserByEmail(player2Email)
	if err != nil {
		log.Errorf("Error fetching user: %v", err)
		return err
	}

	game, err := s.gameRepository.Fetch(key)
	if err != nil {
		return err
	}

	err = game.Join(user)
	if err != nil {
		return err
	}

	s.gameRepository.Save(game)
	return nil
}

func (s GamesService) PlayMove(key string, playerEmail string, column int) error {
	user, err := s.userService.FindUserByEmail(playerEmail)
	if err != nil {
		log.Errorf("Error fetching user: %v", err)
		return err
	}
	game, err := s.gameRepository.Fetch(key)
	if err != nil {
		return err
	}
	err = game.Play(user, column)
	if err != nil {
		return err
	}
	s.gameRepository.Save(game)
	return nil
}

func (s GamesService) NewGame(player1Email string, public bool) NewGameResponse {
	user, err := s.userService.FindUserByEmail(player1Email)
	if err != nil {
		log.Errorf("Error fetching user: %v", err)
		return NewGameResponse{}
	}

	game := model.NewGame(user, public)
	if s.gameRepository.Save(game) {
		return NewGameResponseFromGame(game)
	} else {
		log.Errorf("Error creating new game for player: %s", player1Email)
		return NewGameResponse{}
	}
}
