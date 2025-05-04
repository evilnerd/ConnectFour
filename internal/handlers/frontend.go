package handlers

import (
	"connectfour/internal/service"
	"connectfour/resource/templ"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"time"
)

type Game struct {
	Key       string `json:"key"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	Status    string `json:"status"`
}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Serving landing page (login form)")
	// Set appropriate content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Parse both the main template and the navigation template
	tmpl := template.Must(template.ParseFS(templ.TemplateFS, templ.LoginPage, templ.NavigationBar))
	_ = tmpl.Execute(w, nil)
}

func GamesPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Serving games page")
	// Set appropriate content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Define a template with the custom function
	tmpl := template.New("games.html").Funcs(template.FuncMap{
		"formatDate": formatDate,
	})

	// Parse both the main template and the navigation template
	tmpl = template.Must(tmpl.ParseFS(templ.TemplateFS, templ.GamesPage, templ.NavigationBar))

	// Pass the data to the template
	_ = tmpl.Execute(w, nil)
}

func MyGamesDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Listing all my games")
	email := emailFromContext(r)
	if email == "" {
		errorResponse(w, "No active user", http.StatusUnauthorized)
	}

	// Return templated HTML for HTMX
	games := gamesService.AllMyGames(email)

	// Set the Content-Type to ensure browser treats response as HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	renderGamesList(w, games)
}

func renderGamesList(w http.ResponseWriter, games []service.NewGameResponse) {
	// Convert time.Time to string format for the template
	type Game struct {
		Key       string `json:"key"`
		CreatedAt string `json:"created_at"`
		CreatedBy string `json:"created_by"`
		Status    string `json:"status"`
	}

	var templateGames []Game
	for _, g := range games {
		templateGames = append(templateGames, Game{
			Key:       g.Key,
			CreatedAt: g.CreatedAt.Format(time.RFC3339),
			CreatedBy: g.CreatedBy,
			Status:    string(g.Status),
		})
	}

	// Define a template function to format dates
	funcMap := template.FuncMap{
		"formatDate": func(dateStr string) string {
			date, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				return dateStr
			}
			return date.Format("January 2, 2006 at 3:04 PM")
		},
	}

	// If there are no games, pass nil to the template to trigger the empty state
	var data interface{}
	if len(templateGames) > 0 {
		data = templateGames
	} else {
		data = nil // This will make the "{{if .}}" condition in the template evaluate to false
	}

	// Create template with the function map
	tmpl := template.New(templ.GamesList).Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseFS(templ.TemplateFS, templ.GamesList))

	// Execute the template
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func formatDate(dateStr string) string {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return date.Format("January 2, 2006 at 3:04 PM")
}
