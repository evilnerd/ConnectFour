package backend

import (
	"connectfour/internal/service"
	"errors"
	"fmt"
	"net/http"
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

func InitWebClient() {
	if u, ok := os.LookupEnv("CONNECT_FOUR_SERVER_URL"); ok {
		log.Printf("Using api url override from CONNECT_FOUR_SERVER_URL environment variable: %s\n", u)
		ServerUrl = u
	} else {
		log.Println("No api url override found. Use the environment variable CONNECT_FOUR_SERVER_URL to override url used.")
	}

}
func Hello() error {
	req, err := http.Get(ServerUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error making a request to the api: %v", err))
	}

	if req.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("the api responded with an error: %d - %s", req.StatusCode, req.Status))
	}

	return nil
}

func Login(wc *WebClient, email string, password string) (err error) {
	req := service.LoginRequest{
		Email:    email,
		Password: password,
	}
	return wc.CallLogin(wc.Url("login"), req)
}

// JoinableGames returns a list of games that the player can join.
func JoinableGames(wc *WebClient) []service.NewGameResponse {
	resp := make([]service.NewGameResponse, 0)
	err := wc.Call(
		http.MethodGet,
		wc.Url("games"),
		&resp,
	)
	if err != nil {
		return resp
	}
	return resp
}

// MyGames returns a list of running games that the user is a part of.
func MyGames(wc *WebClient) []service.NewGameResponse {
	resp := make([]service.NewGameResponse, 0)
	err := wc.Call(
		http.MethodGet,
		wc.Url("games", "my"),
		&resp,
	)
	if err != nil {
		return resp
	}
	return resp
}

func CreateGame(wc *WebClient, public bool) service.NewGameResponse {
	var resp service.NewGameResponse
	err := wc.CallWithBody(
		http.MethodPost,
		wc.Url("games"),
		service.NewGameRequest{Public: public},
		&resp,
	)
	if err != nil {
		return service.NewGameResponse{}
	}
	return resp
}

// GameInfo returns the state of the game in a complete struct that contains rich info about the game.
func GameInfo(wc *WebClient, key string) (service.GameStateResponse, error) {
	var resp service.GameStateResponse
	err := wc.Call(
		http.MethodGet,
		wc.Url("games", key),
		&resp,
	)

	if err != nil {
		return service.GameStateResponse{}, err
	}
	return resp, nil
}

// Move plays a move in an existing game.
func Move(wc *WebClient, key string, column int) (service.GameStateResponse, error) {
	req := service.PlayMoveRequest{
		Column: column,
	}
	var resp service.GameStateResponse
	err := wc.CallWithBody(
		http.MethodPost,
		wc.Url("games", key, "play"),
		req,
		&resp,
	)
	if err != nil {
		return service.GameStateResponse{}, err
	}
	return resp, nil
}

// Join tells the api that the player wants to join an existing game.
func Join(wc *WebClient, key string) (service.GameStateResponse, error) {
	var resp service.GameStateResponse
	err := wc.Call(
		http.MethodPost,
		wc.Url("games", key, "join"),
		&resp,
	)
	if err != nil {
		return service.GameStateResponse{}, err
	}
	return resp, nil
}
