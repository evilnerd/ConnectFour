package handlers

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	// Set up static file serving
	SetupStaticFiles(r)

	// Create public routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", LandingPageHandler)
		r.Get("/test", GreetHandler)         // GET  /test
		r.Post("/login", LoginHandler)       // POST /login
		r.Post("/register", RegisterHandler) // POST /login

		// Diagnostic endpoint
		r.Get("/debug-info", DiagnosticsPageHandler)
	})

	r.Route("/ui", func(r chi.Router) {
		// Public UI routes that don't require JWT
		r.Get("/", LandingPageHandler)
		r.Use(JwtValidation)
		r.Get("/games", GamesPageHandler) // GET /ui/games - HTML page
	})

	// Create routes that need authentication, so they check for the jwt token to be there
	r.Route("/games", func(r chi.Router) {
		r.Use(JwtValidation)
		r.Get("/", OpenGamesHandler)           // GET  /games
		r.Get("/my", MyGamesDataHandler)       // GET  /games/my - HTML for HTMX
		r.Get("/data", MyGamesHandler)         // GET  /games/data - JSON data for API
		r.Post("/", NewGameHandler)            // POST /games
		r.Get("/{key}", GameStateHandler)      // GET  /games/1234abcd
		r.Post("/{key}/join", JoinGameHandler) // POST /games/1234abcd/join
		r.Post("/{key}/play", PlayMoveHandler) // POST /games/1234abcd/play
	})
}
