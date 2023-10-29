package main

import (
	. "connectfour/game"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	if err := tea.NewProgram(NewModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)

	}

}
