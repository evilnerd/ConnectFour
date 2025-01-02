package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type AskNameModel struct {
	ConnectFourModel
	*State
	Text textinput.Model
}

func (s Step) String() string { return string(s) }

func NewAskNameModel(state *State) *AskNameModel {
	m := &AskNameModel{
		State: state,
		Text:  textinput.New(),
	}
	m.Text.Placeholder = "Your name"
	m.Text.Focus()
	m.Text.CharLimit = 40
	m.Text.Width = 40
	return m
}

func (m AskNameModel) BreadCrumb() string {
	if m.PlayerName == "" {
		return "Name"
	}
	return fmt.Sprintf("Name (%s)", m.PlayerName)
}

func (m AskNameModel) Init() tea.Cmd {
	return nil
}

func (m AskNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// If a form is active, make that one handle the key updates.
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		case "esc", "ctrl+c":
			return m.PreviousModel()

		case "enter":
			if isValidName(m.Text.Value()) {
				m.PlayerName = m.Text.Value()
				return m.NextModel()
			}
		}
	}

	var cmd tea.Cmd
	m.Text, cmd = m.Text.Update(msg)
	return m, cmd
}

func isValidName(val string) bool {
	l := len(strings.TrimSpace(val))
	return l > 2 && l <= 100
}

func (m AskNameModel) View() string {
	view := lipgloss.JoinVertical(lipgloss.Left,
		styles.Description.Render("Enter your name to uniquely identify you"),
		m.Text.View(),
	)
	return m.CommonView(view)
}
