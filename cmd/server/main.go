package main

import (
	"connectfour/internal/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.Println("ConnectFour Server")

	r := chi.NewRouter()
	handlers.SetupMiddlewares(r)
	handlers.SetupRoutes(r)

	port := port()
	log.Printf("Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Error while running the api: %v", err)
	}
}

func port() string {
	port, ok := os.LookupEnv("CONNECT_FOUR_SERVER_PORT")
	if !ok {
		port = "8443"
	}
	return port
}
