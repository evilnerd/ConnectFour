package main

import (
	handlers2 "connectfour/internal/handlers"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.Println("ConnectFour Server")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// A good base middleware stack
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.Throttle(100))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Create public routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers2.GreetHandler)             // GET /
		r.Post("/login", handlers2.LoginHandler)       // POST /login
		r.Post("/register", handlers2.RegisterHandler) // POST /login
	})

	// Create routes that need authentication, so they check for the jwt token to be there
	r.Route("/games", func(r chi.Router) {
		r.Use(handlers2.JwtValidation)
		r.Get("/", handlers2.AllGamesHandler)            // GET  /games
		r.Post("/", handlers2.NewGameHandler)            // POST /games
		r.Get("/{key}", handlers2.GameStateHandler)      // GET  /games/1234abcd
		r.Post("/{key}/join", handlers2.JoinGameHandler) // POST /games/1234abcd/join
		r.Post("/{key}/play", handlers2.PlayMoveHandler) // POST /games/1234abcd/play
	})

	port, ok := os.LookupEnv("CONNECT_FOUR_SERVER_PORT")
	if !ok {
		port = "8443"
	}
	log.Printf("Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Error while running the api: %v", err)
	}
}
