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

func GameStateHandler(response http.ResponseWriter, request *http.Request) {
	if key, ok := parseAndCheck(response, request); ok {
		marshal(gamesService.GetGameState(key), response)
	}
}

func OpenGamesHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Listing all games")
	email := emailFromContext(request)
	marshal(gamesService.AllOpenGames(email), response)
}

func MyGamesHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Listing all my games")
	email := emailFromContext(request)
	marshal(gamesService.AllMyGames(email), response)
}

func NewGameHandler(response http.ResponseWriter, request *http.Request) {
	if req, ok := unmarshal[service.NewGameRequest](response, request); ok {
		email := emailFromContext(request)
		game := gamesService.NewGame(email, req.Public)
		marshal(game, response)
	}
}

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
