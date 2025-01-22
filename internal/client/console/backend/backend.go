package backend

import (
	"bytes"
	"connectfour/internal/service"
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

// GameNotFoundError is returned when the model, requested by the user doesn't exist.
type GameNotFoundError struct{}

func (g GameNotFoundError) Error() string { return "Game not found" }

// endregion

func CheckSettings() {
	if u, ok := os.LookupEnv("CONNECT_FOUR_SERVER_URL"); ok {
		log.Printf("Using api url override from CONNECT_FOUR_SERVER_URL environment variable: %s\n", u)
		ServerUrl = u
	} else {
		log.Println("No api url override found. Use the environment variable CONNECT_FOUR_SERVER_URL to override url used.")
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
		return errors.New(fmt.Sprintf("there was an error making a request to the api: %v", err))
	}

	if req.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("the api responded with an error: %d - %s", req.StatusCode, req.Status))
	}

	return nil
}

// JoinableGames returns a list of games that the player can join.
func JoinableGames() []service.NewGameResponse {

	req, err := http.Get(makeUrl("model"))
	if err != nil {
		log.Printf("There was an error making a request to the api: %v\n", err)
		return nil
	}

	if req.StatusCode != http.StatusOK {
		log.Printf("The api responded with an error: %d - %s\n", req.StatusCode, req.Status)
		return nil
	}

	dec := json.NewDecoder(req.Body)

	resp := make([]service.NewGameResponse, 0)
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the list of games from the api: %v\n", err)
	}

	return resp
}

func CreateGame(email string, public bool) service.NewGameResponse {
	req := service.NewGameRequest{
		Public: public,
	}

	body, _ := json.Marshal(req)

	response, err := http.Post(makeUrl("model"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("There was an error making a request to the api: %v\n", err)
		return service.NewGameResponse{}
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("The api responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return service.NewGameResponse{}
	}

	dec := json.NewDecoder(response.Body)

	var resp service.NewGameResponse
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the new-model information from the api: %v\n", err)
	}

	return resp

}

// GameInfo returns the state of the model in a complete struct that contains rich info about the model.
func GameInfo(key string) (service.GameStateResponse, error) {
	return gameStateRequest(nil, "model", key)
}

// Move plays a move in an existing model.
func Move(key string, column int) (service.GameStateResponse, error) {
	req := service.PlayMoveRequest{
		Column: column,
	}
	return gameStateRequest(req, "model", key, "play")
}

// Join tells the api that the player wants to join an existing model.
func Join(key string, email string) (service.GameStateResponse, error) {

	return gameStateRequest(nil, "model", key, "join")
}

func gameStateRequest(req any, parts ...string) (service.GameStateResponse, error) {
	var response *http.Response
	var err error

	if req == nil && parts[len(parts)-1] != "join" {
		response, err = http.Get(makeUrl(parts...))
	} else {
		body, _ := json.Marshal(req)
		response, err = http.Post(makeUrl(parts...), "application/json", bytes.NewBuffer(body))
		if err != nil || response == nil {
			log.Printf("There was an error making a request to the api: %v\n", err)
			return service.GameStateResponse{}, err
		}
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("The api responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return service.GameStateResponse{}, err
	}

	dec := json.NewDecoder(response.Body)

	var resp service.GameStateResponse
	err = dec.Decode(&resp)

	if err != nil {
		log.Printf("Could not decode the model-state information from the api: %v\n", err)
		return service.GameStateResponse{}, err
	}

	return resp, nil
}
