package handlers

import (
	"connectfour/server"
	"errors"
	"github.com/Masterminds/goutils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GameState(response http.ResponseWriter, request *http.Request) {
	if key, ok := parseGameKey(response, request); ok {
		marshal(server.GetGameState(key), response)
	}
}

func NewGame(response http.ResponseWriter, request *http.Request) {
	if req, ok := unmarshal[server.NewGameRequest](response, request); ok {
		game := server.NewGame(req.Player1Name, req.Public)
		marshal(server.NewGameResponseFromGame(game), response)
	}
}

func JoinGame(response http.ResponseWriter, request *http.Request) {
	key, ok := parseGameKey(response, request)
	if !ok {
		return
	}

	if req, ok := unmarshal[server.JoinGameRequest](response, request); ok {
		game := server.GetGame(key)
		err := game.Join(req.Player2Name)
		if handleError(err, response) {
			marshal(server.GetGameState(key), response)
		}
	}
}

func PlayMove(response http.ResponseWriter, request *http.Request) {
	key, ok := parseGameKey(response, request)
	if !ok {
		return
	}
	if req, ok := unmarshal[server.PlayMoveRequest](response, request); ok {
		game := server.GetGame(key)
		err := game.Play(req.Column)
		if handleError(err, response) {
			marshal(server.GetGameState(key), response)
		}
	}
}

func parseGameKey(response http.ResponseWriter, request *http.Request) (string, bool) {
	key := chi.URLParam(request, "key")
	if goutils.IsBlank(key) {
		handleError(errors.New("game key parameter is missing from the uri"), response)
		return "", false
	}

	return key, checkGame(key, response)
}

func checkGame(key string, response http.ResponseWriter) bool {
	if server.GameExists(key) {
		return true
	}

	return handleError(server.NewUnknownGameError(key), response)
}
