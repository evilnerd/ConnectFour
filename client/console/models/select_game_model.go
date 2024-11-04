package models

import (
	"connectfour/client/console"
	"connectfour/client/console/backend"
	"connectfour/server"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type GamesFetched struct {
	games []server.NewGameResponse
}

type GameSelected struct {
	GameKey string
}

type SelectGameModel struct {
	Games []server.NewGameResponse
	List  list.Model
}

func NewSelectGameModel() SelectGameModel {
	return SelectGameModel{}
}

func (m SelectGameModel) Init() tea.Cmd {
	return loadGames
}

func (m SelectGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg.(type) {
	case GamesFetched:
		initGamesList(&m)
	case tea.KeyMsg:
		if msg == tea.KeyEnter {
			return m, func() tea.Msg { return GameSelected{GameKey: key} }
		}

	default:
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m SelectGameModel) View() string {
	//TODO implement me
	panic("implement me")
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
