package models

import (
	"connectfour/client/console"
	"connectfour/client/console/backend"
	"connectfour/server"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type GamesFetched struct {
	games []server.NewGameResponse
}

type GameSelected struct {
	GameKey string
}

type SelectGameModel struct {
	*State
	Games   []server.NewGameResponse
	List    list.Model
	loading bool
}

func NewSelectGameModel(state *State) *SelectGameModel {
	return &SelectGameModel{
		State:   state,
		loading: true,
	}
}

func (m SelectGameModel) BreadCrumb() string {
	return "Select"
}
func (m SelectGameModel) Init() tea.Cmd {
	return loadGames
}

func (m SelectGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case GamesFetched:
		m.Games = msg.games
		log.Printf("Games fetched. Size = %d\n", len(msg.games))
		m.loading = false
		initGamesList(&m)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.PreviousModel()
		case "enter":
			if (!m.loading) && len(m.Games) > 0 {
				m.Key = m.List.SelectedItem().(console.Option).Key()
				return m.NextModel()
			} else {
				return m.PreviousModel()
			}
		}
	}

	if !m.loading && len(m.Games) > 0 {
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	} else {
		return m, nil
	}
}

func (m SelectGameModel) View() string {

	contents := ""

	if m.loading {
		contents = "Loading games..."
	} else {
		if len(m.Games) == 0 {
			contents = "There were no open games to choose from."
		} else {
			contents = m.List.View()
		}
	}

	return m.CommonView(lipgloss.JoinVertical(lipgloss.Left,
		styles.Description.Render("Choose a game to (re-)join"),
		contents,
	))
}

func loadGames() tea.Msg {
	games := backend.JoinableGames()
	return GamesFetched{games: games}
}

func initGamesList(m *SelectGameModel) {
	options := make([]list.Item, 0)
	for _, game := range m.Games {
		options = append(options,
			console.NewOption(
				game.Key,
				fmt.Sprintf("%s (%s)", game.CreatedBy, game.Key),
				fmt.Sprintf("Created at %s | status: %s", game.CreatedAt, game.Status)),
		)
	}

	delegate := list.NewDefaultDelegate()

	m.List = list.New(options, delegate, 80, 20)
	m.List.SetShowHelp(false)
	m.List.SetShowStatusBar(false)
	m.List.SetFilteringEnabled(true)
	m.List.SetShowPagination(true)
	m.List.SetShowTitle(false)

}
