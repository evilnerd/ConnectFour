package models

import (
	"connectfour/client/console/backend"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type CreateGameModel struct {
	*State
	created bool
}

func NewCreateGameModel(state *State) *CreateGameModel {
	return &CreateGameModel{
		State:   state,
		created: false,
	}
}

func (m CreateGameModel) Init() tea.Cmd {
	return createGame(m.PlayerName, m.IsPrivateGame)
}

func (m CreateGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case GameCreated:
		m.created = true
		m.Key = msg.game.Key

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			if m.created {
				return m.NextModel()
			}
		}
	}

	return m, nil
}

func (m CreateGameModel) View() string {

	view := styles.Header.Render("Create a new game\n")

	if m.created {
		view += styles.Label.Render("The game has been created, the key is: ") + styles.Value.Render(m.Key)
		view += "\n" + styles.Label.Render("Press enter to continue.")
	} else {
		view += styles.Value.Render("The game is being created...")
	}
	return view
}

func createGame(name string, private bool) tea.Cmd {
	return func() tea.Msg {
		result := backend.CreateGame(name, !private)
		if result.Key == "" {
			// Something went wrong.
			log.Error("Something went wrong creating the game.")
		}
		return GameCreated{game: result}
	}
}
