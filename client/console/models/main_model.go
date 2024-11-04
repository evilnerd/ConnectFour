package models

import (
	"connectfour/server"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	current       tea.Model
	IsNewGame     bool
	IsPrivateGame bool
	Key           string
	GameStatus    server.GameStatus
}

const (
	AskTheName    Step = "get the name"
	StartOrJoin   Step = "start or join"
	SelectGame    Step = "choose an existing game key"
	AskTheGameKey Step = "enter an existing game key"
	ShowGameKey   Step = "show the current game key"
	StartGame     Step = "start the game"
)

type BackStep struct{} // this message is sent when the user pressed ESC.

func NewMainModel() MainModel {
	return MainModel{
		current: NewStartModel(),
	}
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg.(type) {
	case BackStep:
		m.current = m.PreviousModel()
		return m.current, nil
	}

	newModel, cmd := (m.current).Update(msg)
	m.current = newModel
	return newModel, cmd
}

func (m MainModel) View() string {
	return m.current.View()
}

func (m MainModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m MainModel) PreviousModel() tea.Model {
	switch m.current.(type) {
	case AskNameModel:
		return NewStartModel()
	case AskKeyModel:
		return NewStartModel()
	case GameModel:
		return NewGameModel()
	}
	return nil
}
