package main

import (
	"connectfour/client/console/backend"
	. "connectfour/client/console/models"
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

	// Check settings
	backend.CheckSettings()

	log.Println("Creating new program!")
	m := CreateModels()
	p := tea.NewProgram(m)
	p.SetWindowTitle("ConnectFour Online - The Console Edition")

	flag.StringVar(&m.PlayerName, "name", "", "Pre-specify the user's name to skip the first step.")
	flag.StringVar(&m.Key, "key", "", "Pre-specify the game key to join.")
	flag.Parse()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running startup form: ", err)
		os.Exit(1)
	}
}
