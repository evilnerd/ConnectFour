package main

import (
	"connectfour/cmd/server/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"time"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Create routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.Greet)
		r.Post("/game", handlers.NewGame)             // POST /game
		r.Post("/game/{key}/join", handlers.JoinGame) // POST /game/1234abcd
		r.Post("/game/{key}/play", handlers.PlayMove)
		r.Get("/game/{key}", handlers.GameState) // GET /game/1234abcd
		r.Get("/game", handlers.AllGames)        // GET /game
	})

	// Mount the admin sub-router
	//r.Mount("/admin", adminRouter())
	port, ok := os.LookupEnv("CONNECT_FOUR_SERVER_PORT")
	if !ok {
		port = "8443"
	}
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Error while running the server: %v", err)
	}
}
