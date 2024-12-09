package backend

import (
	"bytes"
	"connectfour/server"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"log"
)

var (
	ServerUrl = "http://localhost:8443"
)

// region Custom errors

// GameNotFoundError is returned when the game, requested by the user doesn't exist.
type GameNotFoundError struct{}

func (g GameNotFoundError) Error() string { return "Game not found" }

// endregion

func CheckSettings() {
	if u, ok := os.LookupEnv("CONNECT_FOUR_SERVER_URL"); ok {
		log.Printf("Using server url override from CONNECT_FOUR_SERVER_URL environment variable: %s\n", u)
		ServerUrl = u
	} else {
		log.Println("No server url override found. Use the environment variable CONNECT_FOUR_SERVER_URL to override url used.")
	}
}

func makeUrl(elem ...string) string {
	res, err := url.JoinPath(ServerUrl, elem...)
	if err != nil {
		log.Printf("Error constructing URL (base = %s): %v\n", ServerUrl, err)
		return ""
	}
	return res
}

func Hello() error {
	req, err := http.Get(makeUrl(""))
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error making a request to the server: %v", err))
	}

	if req.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("the server responded with an error: %d - %s", req.StatusCode, req.Status))
	}

	return nil
}

// JoinableGames returns a list of games that the player can join.
func JoinableGames() []server.NewGameResponse {

	req, err := http.Get(makeUrl("game"))
	if err != nil {
		log.Printf("There was an error making a request to the server: %v\n", err)
		return nil
	}

	if req.StatusCode != http.StatusOK {
		log.Printf("The server responded with an error: %d - %s\n", req.StatusCode, req.Status)
		return nil
	}

	dec := json.NewDecoder(req.Body)

	resp := make([]server.NewGameResponse, 0)
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the list of games from the server: %v\n", err)
	}

	return resp
}

// GameState returns the state of the game in a simple String representation.
func GameState(key string) (server.GameStatus, error) {
	info, err := GameInfo(key)

	if err != nil {
		return server.Unknown, err
	}

	return info.Status, nil
}

func CreateGame(name string, public bool) server.NewGameResponse {
	req := server.NewGameRequest{
		Player1Name: name,
		Public:      public,
	}

	body, _ := json.Marshal(req)

	response, err := http.Post(makeUrl("game"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("There was an error making a request to the server: %v\n", err)
		return server.NewGameResponse{}
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("The server responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return server.NewGameResponse{}
	}

	dec := json.NewDecoder(response.Body)

	var resp server.NewGameResponse
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the new-game information from the server: %v\n", err)
	}

	return resp

}

// GameInfo returns the state of the game in a complete struct that contains rich info about the game.
func GameInfo(key string) (server.GameStateResponse, error) {
	return gameStateRequest(nil, "game", key)
}

// Move plays a move in an existing game.
func Move(key string, column int) (server.GameStateResponse, error) {
	req := server.PlayMoveRequest{
		Column: column,
	}
	return gameStateRequest(req, "game", key, "play")
}

// Join tells the server that the player wants to join an existing game.
func Join(key string, name string) (server.GameStateResponse, error) {
	req := server.JoinGameRequest{
		Player2Name: name,
	}
	return gameStateRequest(req, "game", key, "join")
}

func gameStateRequest(req any, parts ...string) (server.GameStateResponse, error) {
	var response *http.Response
	var err error

	if req == nil {
		response, err = http.Get(makeUrl(parts...))
	} else {
		body, _ := json.Marshal(req)
		response, err = http.Post(makeUrl(parts...), "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("There was an error making a request to the server: %v\n", err)
			return server.GameStateResponse{}, err
		}
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("The server responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return server.GameStateResponse{}, err
	}

	dec := json.NewDecoder(response.Body)

	var resp server.GameStateResponse
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the game-state information from the server: %v\n", err)
		return server.GameStateResponse{}, err
	}

	return resp, nil
}
