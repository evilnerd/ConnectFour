package models

import tea "github.com/charmbracelet/bubbletea"

type ExitModel struct {
	*State
	message string
}

func NewExitModel(state *State) *ExitModel {
	return &ExitModel{State: state}
}

func (m ExitModel) BreadCrumb() string {
	return "Exit"
}

func (m ExitModel) Init() tea.Cmd {
	return nil
}

func (m ExitModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ExitModel) View() string {
	view := "Thanks for playing."
	if m.message != "" {
		view = "ConnectFour exited: " + m.message
	}
	return m.CommonView(view)
}
