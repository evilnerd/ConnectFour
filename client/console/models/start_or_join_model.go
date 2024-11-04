package models

import (
	"connectfour/client/console"
	"connectfour/server"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type GameCreated struct {
	game server.NewGameResponse
}

type StartOrJoinModel struct {
	MainModel
	List list.Model
}

func (m StartOrJoinModel) Init() tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func (m StartOrJoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		if msg == tea.KeyEnter {
			m.IsNewGame = m.List.Index() < 2
			m.IsPrivateGame = m.List.Index() == 0 || m.List.Index() == 2

			return m.decideNextModel(), nil
		}
	}
	m.List, cmd = m.List.Update(msg)
	if m.Step == SelectGame && m.Games == nil {
		// get the list of games
		cmd = tea.Batch(cmd, loadGames)
	}
	return m, cmd
}

func (m StartOrJoinModel) View() string {
	view := styles.Label.Render("What kind of game do you want to start? \n\n")
	return view + m.List.View()
}

func (m StartOrJoinModel) decideNextModel() tea.Model {
	if m.IsNewGame {
		if m.IsPrivateGame {
			return NewAskKeyModel()
		} else {
			return nil // NewShowKeyModel()
		}
	} else {
		if m.IsPrivateGame {
			return NewAskKeyModel()
		} else {

		}
		//} else if m.IsPrivateGame {
		//	m.Step = AskTheGameKey
		//} else {
		//	m.Step = SelectGame
		//}
	}
}

func initCreateOrJoinList(m *AskNameModel) {
	options := []list.Item{
		console.NewOption("1. Create new private game", "Creates a new game that will not be public, so you must share the key."),
		console.NewOption("2. Create new public game", "Creates a new game that's going to be listed and open for anyone to join."),
		console.NewOption("3. Join a private game", "Join a game that's not listed, but that you received a key for."),
		console.NewOption("4. Join a public game", "Browse the list of games and join one (this will fetch the list of games)."),
	}

	delegate := list.NewDefaultDelegate()

	m.List = list.New(options, delegate, 80, 20)
	m.List.Title = "Kind of game"
	m.List.SetShowStatusBar(false)
	m.List.SetFilteringEnabled(false)
	m.List.SetShowPagination(false)
}
