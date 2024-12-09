package models

import tea "github.com/charmbracelet/bubbletea"

type ExitModel struct {
	message string
}

func NewExitModel(state *State) *ExitModel {
	return &ExitModel{}
}

func (m ExitModel) Init() tea.Cmd {
	return nil
}

func (m ExitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ExitModel) View() string {
	if m.message != "" {
		return "ConnectFour exited: " + m.message + "\n"
	}
	return "Thanks for playing.\n"
}
