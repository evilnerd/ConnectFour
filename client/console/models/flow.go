package models

import (
	"connectfour/client/console/backend"
	"connectfour/server"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type Step string

type State struct {
	CurrentStep   Step
	CurrentModel  tea.Model
	Key           string
	PlayerName    string
	GameStatus    server.GameStatus
	GameInfo      server.GameStateResponse
	IsNewGame     bool
	IsPrivateGame bool
	Loading       bool
}

const (
	AskTheName    Step = "get the name"
	StartOrJoin   Step = "start or join"
	SelectGame    Step = "choose an existing game key"
	AskTheGameKey Step = "enter an existing game key"
	ShowGameKey   Step = "show the current game key"
	StartGame     Step = "start the game"
)

var (
	state            *State
	askKeyModel      AskKeyModel
	askNameModel     AskNameModel
	createGameModel  CreateGameModel
	exitModel        ExitModel
	playGameModel    PlayGameModel
	mainModel        MainModel
	selectGameModel  SelectGameModel
	startOrJoinModel StartOrJoinModel
)

func CreateModels() *MainModel {

	state = &State{}

	mainModel = *NewMainModel(state)
	askKeyModel = *NewAskKeyModel(state)
	askNameModel = *NewAskNameModel(state)
	createGameModel = *NewCreateGameModel(state)
	exitModel = *NewExitModel(state)
	playGameModel = *NewPlayGameModel(state)
	selectGameModel = *NewSelectGameModel(state)
	startOrJoinModel = *NewStartOrJoinModel(state)

	state.CurrentModel = askNameModel

	return &mainModel
}

func (s *State) CantContinueModel(message string) (tea.Model, tea.Cmd) {
	exitModel.message = message
	return exitModel, tea.Quit
}

func (s *State) PreviousModel() (tea.Model, tea.Cmd) {
	var prevModel tea.Model = askNameModel
	var prevCmd tea.Cmd

	switch (s.CurrentModel).(type) {
	case AskNameModel:
		prevModel = exitModel
		prevCmd = tea.Quit
	case AskKeyModel:
		prevModel = askNameModel
	case PlayGameModel:
		prevModel = startOrJoinModel
	case SelectGameModel:
		prevModel = startOrJoinModel
	}
	log.Printf("[Previous] Current Model = %T, Next Model = %T\n", s.CurrentModel, prevModel)
	s.CurrentModel = prevModel
	return prevModel, prevCmd
}

func (s *State) NextModel() (tea.Model, tea.Cmd) {

	var nextModel tea.Model = askNameModel
	var nextCmd tea.Cmd

	switch (s.CurrentModel).(type) {
	case AskNameModel:
		nextModel = startOrJoinModel
	case AskKeyModel:
		nextModel = playGameModel
	case StartOrJoinModel:
		if s.IsNewGame {
			nextModel = createGameModel
			nextCmd = createGame(s.PlayerName, s.IsPrivateGame)
		} else {
			if s.IsPrivateGame {
				nextModel = askKeyModel
			} else {
				nextModel = selectGameModel
				nextCmd = loadGames
			}
		}
	case CreateGameModel:
		nextModel = playGameModel
		nextCmd = loadGameInfo(s.Key)
	case SelectGameModel:
		log.Printf("Player selected game %s, starting game...\n", s.Key)
		nextModel = playGameModel
		nextCmd = joinGame(s.Key, s.PlayerName)
	}

	log.Printf("[Next] Current Model = %T Next Model = %T\n", s.CurrentModel, nextModel)
	s.CurrentModel = nextModel

	return nextModel, nextCmd
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

type GameInfoMsg struct {
	info         server.GameStateResponse
	errorMessage string
}

// loadGameInfo returns a Cmd that sends a tea.Msg when Game Info is fetched.
func loadGameInfo(key string) tea.Cmd {
	return func() tea.Msg {
		info, err := backend.GameInfo(key)
		msg := GameInfoMsg{
			info: info,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This game key could not be found"
		} else {
			msg.info = info
		}
		return msg
	}
}

// joinGame returns a Cmd that sends a tea.Msg when game is joined.
func joinGame(key string, name string) tea.Cmd {
	return func() tea.Msg {
		info, err := backend.Join(key, name)
		msg := GameInfoMsg{
			info: info,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This game key could not be found"
		} else {
			msg.info = info
		}
		return msg
	}
}
