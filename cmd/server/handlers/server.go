package handlers

import (
	"connectfour/server"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Greet(response http.ResponseWriter, request *http.Request) {
	log.Debug("Sending a greeting")
	response.Write([]byte("Let's play a game."))
}

func AllGames(response http.ResponseWriter, request *http.Request) {
	log.Debug("Listing all games")
	marshal(server.AllOpenGames(), response)
}
