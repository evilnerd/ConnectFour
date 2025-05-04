package handlers

import (
	"connectfour/resource/templ"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// SetupStaticFiles configures static file serving
func SetupStaticFiles(r chi.Router) {
	// Create a file server from the embedded files
	fileServer := http.FileServer(http.FS(templ.TemplateFS))

	// Add a middleware that logs file requests
	fileServerWithLogging := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		fmt.Printf("Static file request: %s\n", r.URL.Path)

		// Forward to the actual file server
		fileServer.ServeHTTP(w, r)
	})

	// Serve JS files from the correct path
	r.Handle("/ui/js/*", http.StripPrefix("/ui/", fileServerWithLogging))
}
