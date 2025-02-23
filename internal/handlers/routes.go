package handlers

import "github.com/go-chi/chi/v5"

func SetupRoutes(r *chi.Mux) {
	// Create public routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", GreetHandler)             // GET /
		r.Post("/login", LoginHandler)       // POST /login
		r.Post("/register", RegisterHandler) // POST /login
	})

	// Create routes that need authentication, so they check for the jwt token to be there
	r.Route("/games", func(r chi.Router) {
		r.Use(JwtValidation)
		r.Get("/", OpenGamesHandler)           // GET  /games
		r.Get("/my", MyGamesHandler)           // GET  /games
		r.Post("/", NewGameHandler)            // POST /games
		r.Get("/{key}", GameStateHandler)      // GET  /games/1234abcd
		r.Post("/{key}/join", JoinGameHandler) // POST /games/1234abcd/join
		r.Post("/{key}/play", PlayMoveHandler) // POST /games/1234abcd/play
	})
}
