package console

import (
	"connectfour/server"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

var (
	styles = NewAppStyles()
)

type StartModel struct {
	LocalPlayerName string
	IsNewGame       bool
	IsPrivateGame   bool
	GameKey         string
	Step            Step
	Games           []server.NewGameResponse
	List            list.Model
	Text            textinput.Model
}

type GamesFetched struct {
	games []server.NewGameResponse
}

type GameCreated struct {
	game server.NewGameResponse
}

type Step string

func (s Step) String() string { return string(s) }

const (
	AskTheName    Step = "get the name"
	StartOrJoin   Step = "start or join"
	SelectGame    Step = "choose an existing game key"
	AskTheGameKey Step = "enter an existing game key"
	ShowGameKey   Step = "show the current game key"
	StartGame     Step = "start the game"
)

func NewStartModel() StartModel {
	m := StartModel{
		Step: AskTheName,
		Text: textinput.New(),
	}
	m.Text.Placeholder = "Your name"
	m.Text.Focus()
	m.Text.CharLimit = 40
	m.Text.Width = 40

	return m
}

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// If a form is active, make that one handle the key updates.
	switch msg := msg.(type) {

	case GamesFetched:
		m.Games = msg.games
		initGamesList(&m)

	case GameCreated:
		if msg.game.Key == "" {
			m.Step = StartOrJoin
		} else {
			m.GameKey = msg.game.Key
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "escape":
			return m, tea.Quit
		case "enter":
			progressStep(&m)
		}

	}

	var cmd tea.Cmd
	switch m.Step {
	case AskTheName, AskTheGameKey:
		m.Text, cmd = m.Text.Update(msg)
		return m, cmd

	case StartOrJoin, SelectGame:
		m.List, cmd = m.List.Update(msg)
		if m.Step == SelectGame && m.Games == nil {
			// get the list of games
			cmd = tea.Batch(cmd, loadGames)
		}
		return m, cmd

	case ShowGameKey:
		if m.GameKey == "" {
			return m, startGame(m.LocalPlayerName, !m.IsPrivateGame)
		} else {
			return m, nil
		}

	case StartGame:
		return m, nil
	}

	return m, nil
}

func startGame(name string, public bool) tea.Cmd {
	return func() tea.Msg {
		result := CreateGame(name, public)
		if result.Key == "" {
			// Something went wrong.
			log.Error("Something went wrong creating the game.")
		}
		return GameCreated{game: result}
	}
}

func loadGames() tea.Msg {
	games := JoinableGames()
	return GamesFetched{games: games}
}

//goland:noinspection GoMixedReceiverTypes
func progressStep(m *StartModel) {
	if m.Step == AskTheName {
		m.LocalPlayerName = m.Text.Value()
		m.Step = StartOrJoin
		initCreateOrJoinList(m)
	} else if m.Step == StartOrJoin {
		m.IsNewGame = m.List.Index() < 2
		m.IsPrivateGame = m.List.Index() == 0 || m.List.Index() == 2
		if m.IsNewGame {
			m.Step = ShowGameKey
		} else if m.IsPrivateGame {
			m.Step = AskTheGameKey
		} else {
			m.Step = SelectGame
		}
	} else if m.GameKey != "" {
		if m.Step == StartOrJoin {
			m.Step = ShowGameKey
		} else {
			m.Step = StartGame
		}
	}
}

func initCreateOrJoinList(m *StartModel) {
	options := []list.Item{
		NewOption("1. Create new private game", "Creates a new game that will not be public, so you must share the key."),
		NewOption("2. Create new public game", "Creates a new game that's going to be listed and open for anyone to join."),
		NewOption("3. Join a private game", "Join a game that's not listed, but that you received a key for."),
		NewOption("4. Join a public game", "Browse the list of games and join one (this will fetch the list of games)."),
	}

	delegate := list.NewDefaultDelegate()

	m.List = list.New(options, delegate, 80, 20)
	m.List.Title = "Kind of game"
	m.List.SetShowStatusBar(false)
	m.List.SetFilteringEnabled(false)
	m.List.SetShowPagination(false)
}

func initGamesList(m *StartModel) {
	options := []list.Item{}
	for _, game := range m.Games {
		options = append(options,
			NewOption(
				fmt.Sprintf("%s (%s)", game.CreatedBy, game.Key),
				fmt.Sprintf("Created at %s | status: %s", game.CreatedAt, game.Status)),
		)
	}

	delegate := list.NewDefaultDelegate()

	m.List = list.New(options, delegate, 80, 20)
	m.List.Title = "Select a game to join"
	m.List.SetShowStatusBar(true)
	m.List.SetFilteringEnabled(true)
	m.List.SetShowPagination(true)
}

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) View() string {

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
	default:
		b.WriteString("step = " + m.Step.String())
	}

	return b.String()
	//return lipgloss.JoinVertical(lipgloss.Left, title, form)
}

func (m StartModel) ViewShowGameKey() string {

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
func (m StartModel) ViewAskName() string {
	view := styles.Label.Render("Enter your name")
	return view + "\n" + m.Text.View()
}

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) ViewStartOrJoin() string {
	view := styles.Label.Render("What kind of game do you want to start? \n\n")
	return view + m.List.View()
}

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) ViewSelectGame() string {

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

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) Init() tea.Cmd {
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (m *StartModel) nextStep() {
	if m.Step == AskTheName {
		m.LocalPlayerName = "Dick"
		m.Step = StartOrJoin
	}
}
