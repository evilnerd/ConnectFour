package templ

import (
	"embed"
	"io/fs"
)

// Embed all files including JS
//
// Make sure js files are in the right place
//
//go:embed *.html js/
var templateFS embed.FS

// Public FS is what we expose for direct file access
var TemplateFS, _ = fs.Sub(templateFS, ".")

// Constants for template names
const (
	LoginPage     = "login.html"
	GamesPage     = "games.html"
	GamesList     = "games_list.html"
	NavigationBar = "navigation.html"
)
