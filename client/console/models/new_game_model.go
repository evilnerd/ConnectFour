package models

import (
	"connectfour/client/console/backend"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type NewGameModel struct {
	MainModel
}

func (m NewGameModel) Init() tea.Cmd {
	return loadGames
}

func (m NewGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

func (m NewGameModel) View() string {
	//TODO implement me
	panic("implement me")
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
