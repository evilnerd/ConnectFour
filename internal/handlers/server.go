package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GreetHandler godoc
// @Summary Server greeting endpoint
// @Description Returns a greeting message to confirm the server is running
// @Tags general
// @Accept json
// @Produce text/plain
// @Success 200 {string} string "Let's play a game."
// @Router / [get]
func GreetHandler(response http.ResponseWriter, request *http.Request) {
	log.Debug("Sending a greeting")
	response.Write([]byte("Let's play a game."))
}
