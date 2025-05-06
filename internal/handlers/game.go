package handlers

import (
	"connectfour/internal/model"
	"connectfour/internal/service"
	"errors"
	"github.com/Masterminds/goutils"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GameStateHandler godoc
// @Summary Get game state
// @Description Returns the current state of a specific game
// @Tags games
// @Accept json
// @Produce json
// @Param key path string true "Game key"
// @Security BearerAuth
// @Success 200 {object} service.GameStateResponse "Game state"
// @Failure 400 {object} service.ErrorResponse "Invalid game key"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Failure 404 {object} service.ErrorResponse "Game not found"
// @Router /games/{key} [get]
func GameStateHandler(response http.ResponseWriter, request *http.Request) {
	if key, ok := parseAndCheck(response, request); ok {
		marshal(gamesService.GetGameState(key), response)
	}
}

// OpenGamesHandler godoc
// @Summary List open games
// @Description Returns a list of public games that are open to join
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} service.GameStateResponse "List of open games"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Router /games [get]
func OpenGamesHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Listing all games")
	email := emailFromContext(request)
	marshal(gamesService.AllOpenGames(email), response)
}

// MyGamesHandler godoc
// @Summary List my games
// @Description Returns a list of games where the authenticated user is a player
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} service.GameStateResponse "List of user's games"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Router /games/my [get]
func MyGamesHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Listing all my games")
	email := emailFromContext(request)
	marshal(gamesService.AllMyGames(email), response)
}

// NewGameHandler godoc
// @Summary Create new game
// @Description Creates a new game with the authenticated user as player 1
// @Tags games
// @Accept json
// @Produce json
// @Param game body service.NewGameRequest true "Game options"
// @Security BearerAuth
// @Success 200 {object} service.GameStateResponse "Game created"
// @Failure 400 {object} service.ErrorResponse "Invalid request format"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Router /games [post]
func NewGameHandler(response http.ResponseWriter, request *http.Request) {
	if req, ok := unmarshal[service.NewGameRequest](response, request); ok {
		email := emailFromContext(request)
		game := gamesService.NewGame(email, req.Public)
		marshal(game, response)
	}
}

// JoinGameHandler godoc
// @Summary Join a game
// @Description Allows the authenticated user to join a game as player 2
// @Tags games
// @Accept json
// @Produce json
// @Param key path string true "Game key"
// @Security BearerAuth
// @Success 200 {object} service.GameStateResponse "Successfully joined the game"
// @Failure 400 {object} service.ErrorResponse "Game is already full or not in a joinable state"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Failure 404 {object} service.ErrorResponse "Game not found"
// @Router /games/{key}/join [post]
func JoinGameHandler(response http.ResponseWriter, request *http.Request) {
	key, ok := parseAndCheck(response, request)
	if !ok {
		return
	}

	email := emailFromContext(request)
	err := gamesService.JoinGame(key, email)
	if handleError(err, response) {
		marshal(gamesService.GetGameState(key), response)
	}
}

// PlayMoveHandler godoc
// @Summary Play a move
// @Description Allows the current player to make a move in a game the player is part of.
// @Tags games
// @Accept json
// @Produce json
// @Param key path string true "Game key"
// @Param move body service.PlayMoveRequest true "Move details"
// @Security BearerAuth
// @Success 200 {object} service.GameStateResponse "Move successfully played"
// @Failure 400 {object} service.ErrorResponse "Invalid move or not player's turn"
// @Failure 401 {object} service.ErrorResponse "Unauthorized"
// @Failure 404 {object} service.ErrorResponse "Game not found"
// @Router /games/{key}/play [post]
func PlayMoveHandler(response http.ResponseWriter, request *http.Request) {
	key := parseGameKey(response, request)
	if !checkGame(key, response) {
		return
	}
	if req, ok := unmarshal[service.PlayMoveRequest](response, request); ok {
		email := emailFromContext(request)
		err := gamesService.PlayMove(key, email, req.Column)
		if handleError(err, response) {
			marshal(gamesService.GetGameState(key), response)
		}
	}
}

func parseGameKey(response http.ResponseWriter, request *http.Request) string {
	key := chi.URLParam(request, "key")
	if goutils.IsBlank(key) {
		handleError(errors.New("game key parameter is missing from the uri"), response)
		return ""
	}
	return key
}

func checkGame(key string, response http.ResponseWriter) bool {
	if gamesService.GameExists(key) {
		return true
	}

	return handleError(model.NewUnknownGameError(key), response)
}

func parseAndCheck(response http.ResponseWriter, request *http.Request) (string, bool) {
	key := parseGameKey(response, request)
	return key, checkGame(key, response)
}
