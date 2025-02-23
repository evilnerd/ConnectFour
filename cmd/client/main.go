package main

import (
	"connectfour/internal/client/console/backend"
	. "connectfour/internal/client/console/models"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

func main() {

	file, err := tea.LogToFile("./connectfour.log", "cf")
	if err != nil {
		panic(fmt.Sprintf("Can't open log output: %v\n", err))
	}
	defer file.Close()

	var key string
	var enableJwtFileStorage bool
	//	flag.StringVar(&m.PlayerName, "credentials", "", "Pre-specify the user's credentials to skip the first step (use the form email:password, e.g. --credentials player@email.com:thepassword).")
	flag.StringVar(&key, "key", "", "Pre-specify the game key to join.")
	flag.BoolVar(&enableJwtFileStorage, "store-credentials", true, "Enable jwt file storage (disable to force login screen)")
	flag.Parse()

	// Check settings
	backend.InitWebClient()

	log.Println("Creating new program!")
	m := CreateModels(key, enableJwtFileStorage)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running startup form: ", err)
		os.Exit(1)
	}
}
