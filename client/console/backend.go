package client

import (
	"connectfour/server"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	ServerUrl = "http://localhost:8443"
)

func init() {
	if url, ok := os.LookupEnv("CONNECT_FOUR_SERVER_URL"); ok {
		ServerUrl = url
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
