package client

import (
	"connectfour/server"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
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

type Item struct {
	title       string
	description string
}

func NewItem(title string, description string) Item {
	return Item{
		title, description,
	}
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	return i.description
}

func (i Item) FilterValue() string {
	return i.title
}

type Step string

func (s Step) String() string { return string(s) }

const (
	GetTheName    Step = "get the name"
	StartOrJoin   Step = "start or join"
	SelectGame    Step = "choose an existing game key"
	GetTheGameKey Step = "enter an existing game key"
	ShowGameKey   Step = "show the current game key"
	StartGame     Step = "start the game"
)

func NewStartModel() StartModel {
	m := StartModel{
		Step: GetTheName,
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

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "escape":
			return m, tea.Quit
			return m, tea.Quit
		case "enter":
			progressStep(&m)

		}

	}

	var cmd tea.Cmd
	switch m.Step {
	case GetTheName, GetTheGameKey:
		m.Text, cmd = m.Text.Update(msg)
		return m, cmd

	case StartOrJoin, SelectGame:
		m.List, cmd = m.List.Update(msg)
		if m.Step == SelectGame && m.Games == nil {
			// get the list of games
			cmd = tea.Batch(cmd, loadGames)
		}
		return m, cmd
	}

	return m, nil
}

func loadGames() tea.Msg {
	time.Sleep(time.Second * 2)
	//m.Games = JoinableGames()
	games := []server.NewGameResponse{
		{
			Key:       "AAAAAA",
			CreatedAt: time.Now().Add(-time.Hour * 2),
			CreatedBy: "Dick",
			Status:    server.Created,
		},
		{
			Key:       "BBBBB",
			CreatedAt: time.Now().Add(-time.Hour * 3),
			CreatedBy: "Marian",
			Status:    server.Created,
		},
		{
			Key:       "CCCCC",
			CreatedAt: time.Now().Add(-time.Hour * 4),
			CreatedBy: "Sanae",
			Status:    server.Created,
		},
		{
			Key:       "DDDDDD",
			CreatedAt: time.Now().Add(-time.Hour * 5),
			CreatedBy: "Lucy",
			Status:    server.Created,
		},
	}
	return GamesFetched{games: games}
}

//goland:noinspection GoMixedReceiverTypes
func progressStep(m *StartModel) {
	if m.Step == GetTheName {
		m.LocalPlayerName = m.Text.Value()
		m.Step = StartOrJoin
		initCreateOrJoinList(m)
	} else if m.Step == StartOrJoin {
		m.IsNewGame = m.List.Index() < 2
		m.IsPrivateGame = m.List.Index() == 0 || m.List.Index() == 2
		if m.IsNewGame {
			m.Step = ShowGameKey
		} else if m.IsPrivateGame {
			m.Step = GetTheGameKey
		} else {
			m.Step = SelectGame
		}
	}
}

func initCreateOrJoinList(m *StartModel) {
	options := []list.Item{
		NewItem("1. Create new private game", "Creates a new game that will not be public, so you must share the key."),
		NewItem("2. Create new public game", "Creates a new game that's going to be listed and open for anyone to join."),
		NewItem("3. Join a private game", "Join a game that's not listed, but that you received a key for."),
		NewItem("4. Join a public game", "Browse the list of games and join one (this will fetch the list of games)."),
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
			NewItem(
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
	if m.Step != GetTheName {
		b.WriteRune('\n')
		b.WriteString(styles.Label.Render("Player name"))
		b.WriteString(styles.Value.Render(m.LocalPlayerName))
		b.WriteRune('\n')
	}

	// Step-specific views
	switch m.Step {
	case GetTheName:
		b.WriteString(m.ViewGetName())

	case StartOrJoin:
		b.WriteString(m.ViewStartOrJoin())
	case SelectGame:
		b.WriteString(m.ViewSelectGame())
	default:
		b.WriteString("step = " + m.Step.String())
	}

	return b.String()
	//return lipgloss.JoinVertical(lipgloss.Left, title, form)
}

//goland:noinspection GoMixedReceiverTypes
func (m StartModel) ViewGetName() string {
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
	if m.Step == GetTheName {
		m.LocalPlayerName = "Dick"
		m.Step = StartOrJoin
	}
}
