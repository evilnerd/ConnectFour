package handlers

import (
	"connectfour/server"
	"net/http"
)

func Greet(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Let's play a game."))
}

func AllGames(response http.ResponseWriter, request *http.Request) {
	marshal(server.AllOpenGames(), response)
}
