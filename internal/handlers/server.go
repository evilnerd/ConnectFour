package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GreetHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Sending a greeting")
	response.Write([]byte("Let's play a game."))
}
