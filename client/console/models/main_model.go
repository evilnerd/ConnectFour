package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

type MainModel struct {
	*State
}

func NewMainModel(state *State) *MainModel {
	m := &MainModel{
		State: state,
	}
	return m
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg.(type) {
	case BackStepMsg:
		return m.PreviousModel()

	case QuitMsg:
		return m, tea.Quit
	}

	return m.CurrentModel.Update(msg)
}

func (m MainModel) View() string {
	var b strings.Builder

	// Always write the header.
	b.WriteString(styles.AppTitle.Render("ConnectFour 0.1"))
	b.WriteRune('\n')

	// Write the contents for the current step.
	b.WriteString(m.CurrentModel.View())
	return b.String()
}

func (m MainModel) Init() tea.Cmd {
	log.Printf("Init for model %T\n", m.CurrentModel)
	return m.CurrentModel.Init()
}
