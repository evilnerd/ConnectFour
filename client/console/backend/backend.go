package backend

import (
	"bytes"
	"connectfour/server"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	ServerUrl = "http://localhost:8443"
)

// region Custom errors
type GameNotFoundError struct{}

func (g GameNotFoundError) Error() string { return "Game not found" }

// endregion

func init() {
	if url, ok := os.LookupEnv("CONNECT_FOUR_SERVER_URL"); ok {
		ServerUrl = url
	}
}

func makeUrl(elem ...string) string {
	res, err := url.JoinPath(ServerUrl, elem...)
	if err != nil {
		log.Errorf("Error constructing URL (base = %s): %v\n", ServerUrl, err)
		return ""
	}
	return res
}

func JoinableGames() []server.NewGameResponse {

	req, err := http.Get(makeUrl("game"))
	if err != nil {
		log.Errorf("There was an error making a request to the server: %v\n", err)
		return nil
	}

	if req.StatusCode != http.StatusOK {
		log.Errorf("The server responded with an error: %d - %s\n", req.StatusCode, req.Status)
		return nil
	}

	dec := json.NewDecoder(req.Body)

	resp := make([]server.NewGameResponse, 0)
	err = dec.Decode(&resp)

	if err != nil {
		log.Errorf("Could not decode the list of games from the server: %v\n", err)
	}

	return resp
}

func GameState(key string) (server.GameStatus, error) {
	req, err := http.Get(makeUrl("game", key))
	if err != nil {
		log.Errorf("There was an error making a request to the server: %v\n", err)
		return server.Unknown, err
	}

	if req.StatusCode != http.StatusOK {
		if req.StatusCode == http.StatusNotFound || req.StatusCode == http.StatusBadRequest {
			return server.Unknown, GameNotFoundError{}
		}
		log.Errorf("The server responded with an error: %d - %s\n", req.StatusCode, req.Status)
		return server.Unknown, err
	}

	dec := json.NewDecoder(req.Body)
	var state server.GameStateResponse
	err = dec.Decode(&state)

	if err != nil {
		log.Errorf("Could not decode the status of the game: %v\n", err)
		return server.Unknown, err
	}

	return state.Status, nil
}

func CreateGame(name string, public bool) server.NewGameResponse {
	req := server.NewGameRequest{
		Player1Name: name,
		Public:      public,
	}

	body, _ := json.Marshal(req)

	response, err := http.Post(makeUrl("game"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("There was an error making a request to the server: %v\n", err)
		return server.NewGameResponse{}
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("The server responded with an error: %d - %s\n", response.StatusCode, response.Status)
		return server.NewGameResponse{}
	}

	dec := json.NewDecoder(response.Body)

	var resp server.NewGameResponse
	err = dec.Decode(&resp)

	if err != nil {
		log.Errorf("Could not decode the new-game information from the server: %v\n", err)
	}

	return resp

}
