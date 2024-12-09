package main

import (
	"connectfour/client/console/backend"
	. "connectfour/client/console/models"
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
	p := tea.NewProgram(CreateModels())
	p.SetWindowTitle("ConnectFour Online - The Console Edition")

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running startup form: ", err)
		os.Exit(1)
	}
}
