package models

import (
	"connectfour/internal/client/console"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StartOrJoinModel struct {
	*State
	List list.Model
}

func (m StartOrJoinModel) BreadCrumb() string {
	return "Type"
}

func NewStartOrJoinModel(state *State) *StartOrJoinModel {
	options := []list.Item{
		console.NewOption("1", "1. List my active games", "Lists the games you are participating in, for you to continue."),
		console.NewOption("2", "2. Create new private game", "Creates a new game that will not be public, so you must share the key."),
		console.NewOption("3", "3. Create new public game", "Creates a new game that's going to be listed and open for anyone to join."),
		console.NewOption("4", "4. Join a private game", "Join a game that's not listed, but that you received a key for."),
		console.NewOption("5", "5. Join a public game", "Browse the list of games and join one (this will fetch the list of games)."),
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(options, delegate, 120, 18)
	l.Title = "Kind of game"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	return &StartOrJoinModel{
		List:  l,
		State: state,
	}
}

func (m StartOrJoinModel) Init() tea.Cmd {
	return nil
}

func (m StartOrJoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	i := m.List.Index()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m.PreviousModel()

		case "enter":
			m.IsContinue = i == 0
			m.IsNewGame = i == 1 || i == 2
			m.IsPrivateGame = i == 1 || i == 3

			return m.NextModel()
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m StartOrJoinModel) View() string {
	view := lipgloss.JoinVertical(lipgloss.Left,
		styles.Description.Render("What would you like to do?"),
		m.List.View(),
	)
	return m.CommonView(view)
}
