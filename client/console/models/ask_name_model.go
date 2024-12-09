package models

import (
	"connectfour/client/console"
	"connectfour/client/console/backend"
	"connectfour/server"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	styles = console.NewAppStyles()
)

type AskNameModel struct {
	*State
	Text             textinput.Model
	ConnectionTested bool
	ConnectionError  string
}

func (s Step) String() string { return string(s) }

func NewAskNameModel(state *State) *AskNameModel {
	m := &AskNameModel{
		State:            state,
		Text:             textinput.New(),
		ConnectionTested: false,
		ConnectionError:  "",
	}
	m.Text.Placeholder = "Your name"
	m.Text.Focus()
	m.Text.CharLimit = 40
	m.Text.Width = 40
	return m
}

type NotConnected struct {
	message string
}

type Connected struct{}

func (m AskNameModel) Init() tea.Cmd {
	return connect()
}

func connect() tea.Cmd {
	return func() tea.Msg {
		err := backend.Hello()
		if err != nil {
			return NotConnected{err.Error()}
		}
		return Connected{}
	}
}

func (m AskNameModel) isConnected() bool {
	return m.ConnectionTested && m.ConnectionError == ""
}

func (m AskNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// If a form is active, make that one handle the key updates.
	switch msg := msg.(type) {

	case NotConnected:
		m.ConnectionError = msg.message
		m.ConnectionTested = true

	case Connected:
		m.ConnectionTested = true

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		case "escape", "ctrl+c":
			return m.PreviousModel()

		case "enter":
			m.PlayerName = m.Text.Value()
			return m.NextModel()
		}
	}
	if m.isConnected() {
		var cmd tea.Cmd
		m.Text, cmd = m.Text.Update(msg)
		return m, cmd
	} else if m.ConnectionTested && m.ConnectionError != "" {
		return m, tea.Quit
	}
	return m, nil
}

type GameStateMsg struct {
	tea.Msg
	status       server.GameStatus
	errorMessage string
}

func (m AskNameModel) View() string {
	view := ""
	if !m.ConnectionTested {
		view = styles.Label.Render("Checking connection...\n")
	} else if m.ConnectionError != "" {
		view = styles.Label.Render("Connection error: " + m.ConnectionError + "\n")
	} else {
		view = styles.Label.Render("Enter your name")
		view += "\n" + m.Text.View()
	}
	return view
}
