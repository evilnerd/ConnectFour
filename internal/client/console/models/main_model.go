package models

import (
	"connectfour/internal/client/console"
	"connectfour/internal/client/console/backend"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

var (
	styles = console.NewAppStyles()
)

type MainModel struct {
	*State
	ConnectionTested bool
	ConnectionError  string
}

func NewMainModel(state *State) *MainModel {
	m := &MainModel{
		State:            state,
		ConnectionTested: false,
		ConnectionError:  "",
	}
	return m
}

func (m MainModel) BreadCrumb() string {
	return "Main"
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case NotConnected:
		m.ConnectionError = msg.message
		m.ConnectionTested = true
		return m, QuitCmd

	case Connected:
		m.ConnectionTested = true

	case BackStepMsg:
		return m.PreviousModel()

	case QuitMsg:
		return m, tea.Quit
	}

	if m.ConnectionTested && m.ConnectionError == "" {
		return m.NextModel()
	}

	return m, nil
}

func (m MainModel) View() string {
	var b strings.Builder

	if !m.hasValidFlags() {
		b.WriteString("If you specify the key of a model to join, then you must also specify the player's name.")
		b.WriteRune('\n')
	} else if !m.ConnectionTested {
		b.WriteString(styles.Label.Render("Checking connection..."))
		b.WriteRune('\n')
	} else if m.ConnectionError != "" {
		b.WriteString(styles.Label.Render("Connection error: " + m.ConnectionError))
		b.WriteRune('\n')
	}

	return m.CommonView(b.String())
}

func (m MainModel) Init() tea.Cmd {
	cmd := Connect()

	if m.PlayerName != "" {
		m.SkipAskName()
	}

	if m.Key != "" {
		if m.PlayerName == "" {
			tea.Printf("Problem: If you specify the key, then you must also specify the player's name.")
			m.CurrentModel = exitModel
			return tea.Quit
		}
		m.SkipAskKey()
		cmd = joinGame(m.Key, m.PlayerName)
	}

	log.Printf("Init for model %T\n", m.CurrentModel)
	return cmd
}

func (m MainModel) isConnected() bool {
	return m.ConnectionTested && m.ConnectionError == ""
}

func (m MainModel) hasValidFlags() bool {
	return m.Key == "" || m.PlayerName != ""
}

func Connect() tea.Cmd {
	return func() tea.Msg {
		err := backend.Hello()
		if err != nil {
			return NotConnected{err.Error()}
		}
		return Connected{}
	}
}
