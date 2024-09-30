package main

import (
	. "connectfour/client"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	p := tea.NewProgram(NewStartModel())
	p.SetWindowTitle("ConnectFour Online - The Console Edition")

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running startup form: ", err)
		os.Exit(1)
	}

	//if err := tea.NewProgram(NewGameModel()).Start(); err != nil {
	//	fmt.Println("Error running program:", err)
	//	os.Exit(1)
	//}

}
