package models

import (
	"connectfour/server"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AskKeyModel struct {
	*State
	InputGameKey string
	ErrorMessage string
	Text         textinput.Model // the input box
}

func NewAskKeyModel(state *State) *AskKeyModel {
	m := &AskKeyModel{
		State: state,
		Text:  textinput.New(),
	}
	m.Text.Placeholder = "Your name"
	m.Text.Focus()
	m.Text.CharLimit = 14
	m.Text.Width = 20
	return m
}

func (m AskKeyModel) BreadCrumb() string {
	return "Key"
}

func (m AskKeyModel) Init() tea.Cmd {
	return nil
}

func (m AskKeyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m.PreviousModel()
		case "enter":
			m.Key = m.Text.Value()
			return m.NextModel()
		}
	}
	m.Text, _ = m.Text.Update(msg)
	return m, nil
}

func (m AskKeyModel) View() string {

	view := ""
	if m.InputGameKey == "" && m.ErrorMessage == "" {
		view = lipgloss.JoinVertical(lipgloss.Left,
			styles.Label.Render("Enter the private game key that was shared with you"),
			m.Text.View(),
		)
	} else {
		if m.ErrorMessage != "" {
			view = styles.Label.Render(m.ErrorMessage) + "\n"
		} else if m.GameStatus == server.Unknown {
			view = styles.Subdued.Render("Looking for that game...\n")
		} else {
			view = styles.Subdued.Render("The game status is " + string(m.GameStatus) + "\n")
		}
		view = lipgloss.JoinVertical(lipgloss.Left,
			styles.Label.Render("Game key ")+styles.Value.Render(m.InputGameKey),
			view,
		)
	}

	return m.CommonView(view)
}
