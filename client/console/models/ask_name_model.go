package models

import (
	"connectfour/client/console"
	"connectfour/client/console/backend"
	"connectfour/server"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

var (
	styles = console.NewAppStyles()
)

type AskNameModel struct {
	LocalPlayerName string
	IsNewGame       bool
	IsPrivateGame   bool
	GameKey         string
	GameStatus      server.GameStatus
	ErrorMessage    string
	Step            Step
	Games           []server.NewGameResponse
	List            list.Model
	Text            textinput.Model
}

type Step string

func (s Step) String() string { return string(s) }

func NewStartModel() AskNameModel {
	m := AskNameModel{
		Step: AskTheName,
		Text: textinput.New(),
	}
	m.Text.Placeholder = "Your name"
	m.Text.Focus()
	m.Text.CharLimit = 40
	m.Text.Width = 40
	m.resetGame()
	return m
}

//goland:noinspection GoMixedReceiverTypes
func (m AskNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// If a form is active, make that one handle the key updates.
	switch msg := msg.(type) {

	case GameCreated:
		if msg.game.Key == "" {
			m.Step = StartOrJoin
		} else {
			m.GameKey = msg.game.Key
		}

	case GameStateMsg:
		m.GameStatus = msg.status
		m.ErrorMessage = msg.errorMessage
		if msg.errorMessage != "" {
			m.GameKey = ""
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "escape":
			return m, tea.Quit
		case "enter":
			m.progressStep()
		}
	}

	var cmd tea.Cmd
	switch m.Step {
	case AskTheName:
		m.Text, cmd = m.Text.Update(msg)
		return m, cmd

	case AskTheGameKey:
		if m.GameKey == "" {
			m.Text, cmd = m.Text.Update(msg)
			return m, cmd
		} else {
			return m, loadGameState(m.GameKey)
		}

	case ShowGameKey:
		if m.GameKey == "" {
			return m, createGame(m.LocalPlayerName, !m.IsPrivateGame)
		} else {
			return m, nil
		}

	case StartGame:
		return m, nil
	}

	return m, nil
}

func createGame(name string, public bool) tea.Cmd {
	return func() tea.Msg {
		result := backend.CreateGame(name, public)
		if result.Key == "" {
			// Something went wrong.
			log.Error("Something went wrong creating the game.")
		}
		return GameCreated{game: result}
	}
}

type GameStateMsg struct {
	tea.Msg
	status       server.GameStatus
	errorMessage string
}

func loadGameState(key string) tea.Cmd {
	return func() tea.Msg {
		s, err := backend.GameState(key)
		msg := GameStateMsg{
			status: server.Unknown,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This game key could not be found"
		} else {
			msg.status = s
		}
		return msg
	}
}

//goland:noinspection GoMixedReceiverTypes
func (m *AskNameModel) progressStep() {
	switch m.Step {
	case AskTheName:
		m.LocalPlayerName = m.Text.Value()
		m.Text.SetValue("")
		m.Step = StartOrJoin
		initCreateOrJoinList(m)
		break

	case StartOrJoin:
		m.IsNewGame = m.List.Index() < 2
		m.IsPrivateGame = m.List.Index() == 0 || m.List.Index() == 2
		m.resetGame()
		if m.IsNewGame {
			m.Step = ShowGameKey
		} else if m.IsPrivateGame {
			m.Step = AskTheGameKey
		} else {
			m.Step = SelectGame
		}
		break

	case AskTheGameKey:
		m.ErrorMessage = ""
		if m.GameKey == "" {
			m.GameKey = strings.TrimSpace(m.Text.Value())
			m.Text.Reset()
		}
		if m.GameKey != "" && m.GameStatus == server.Unknown {
			s, err := backend.GameState(m.GameKey)
			if errors.Is(err, backend.GameNotFoundError{}) {
				m.resetGame()
				m.ErrorMessage = "This game key could not be found"
			} else {
				m.GameStatus = s
			}
		}
	}
}

func (m *AskNameModel) resetGame() {
	m.GameStatus = server.Unknown
	m.GameKey = ""
}

//goland:noinspection GoMixedReceiverTypes
func (m AskNameModel) View() string {

	var b strings.Builder

	// Title
	b.WriteString(styles.AppTitle.Render("ConnectFour 0.1"))

	// Player name
	if m.Step != AskTheName {
		b.WriteRune('\n')
		b.WriteString(styles.Label.Render("Player name"))
		b.WriteString(styles.Value.Render(m.LocalPlayerName))
		b.WriteRune('\n')
	}

	// Step-specific views
	switch m.Step {
	case AskTheName:
		b.WriteString(m.ViewAskName())

	case StartOrJoin:
		b.WriteString(m.ViewStartOrJoin())
	case SelectGame:
		b.WriteString(m.ViewSelectGame())
	case ShowGameKey:
		b.WriteString(m.ViewShowGameKey())

	case AskTheGameKey:
		b.WriteString(m.ViewAskGameKey())
	default:
		b.WriteString("step = " + m.Step.String())
	}

	return b.String()
	//return lipgloss.JoinVertical(lipgloss.Left, title, form)
}

// region Specific views
func (m AskNameModel) ViewShowGameKey() string {

	view := styles.Label.Render("Your game key is ")

	if m.GameKey == "" {
		view += styles.Value.Render("being generated")
	} else {
		view += styles.Value.Render(m.GameKey)
	}

	// TODO: if the game is not yet used by a player 2 we can't start the game yet.

	return view
}

//goland:noinspection GoMixedReceiverTypes
func (m AskNameModel) ViewAskName() string {
	view := styles.Label.Render("Enter your name")
	return view + "\n" + m.Text.View()
}

//goland:noinspection GoMixedReceiverTypes
func (m AskNameModel) ViewSelectGame() string {

	var b strings.Builder

	if m.Games == nil {
		b.WriteString("Loading games...")
	} else {
		if len(m.Games) == 0 {
			b.WriteString("There were no open games to choose from.")
		} else {
			b.WriteString(m.List.View())
		}
	}

	return b.String()
}

func (m AskNameModel) ViewAskGameKey() string {
	view := ""
	if m.GameKey == "" && m.ErrorMessage == "" {
		view = styles.Label.Render("Enter the private game key that was shared with you")
		view += "\n" + m.Text.View()
	} else {
		view += styles.Label.Render("Game key ") + styles.Value.Render(m.GameKey) + "\n"
		if m.ErrorMessage != "" {
			view += styles.Label.Render(m.ErrorMessage) + "\n"
		} else if m.GameStatus == server.Unknown {
			view += styles.Subdued.Render("Looking for that game...\n")
		} else {
			view += styles.Subdued.Render("The game status is " + string(m.GameStatus) + "\n")
		}
	}

	return view
}

// endregion

//goland:noinspection GoMixedReceiverTypes
func (m AskNameModel) Init() tea.Cmd {
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (m *AskNameModel) nextStep() {
	if m.Step == AskTheName {
		m.LocalPlayerName = "Dick"
		m.Step = StartOrJoin
	}
}
