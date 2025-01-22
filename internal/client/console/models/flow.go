package models

import (
	"connectfour/internal/client/console/backend"
	"connectfour/internal/model"
	"connectfour/internal/service"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type BreadCrumber interface {
	BreadCrumb() string
}

type ConnectFourModel interface {
	tea.Model
	BreadCrumber
}

type StepNode struct {
	BreadCrumb string
	Step       Step
	Previous   *StepNode
	Model      *ConnectFourModel
}

type Step string

type State struct {
	CurrentStep   *StepNode
	CurrentModel  ConnectFourModel
	Key           string
	PlayerName    string
	GameStatus    model.GameStatus
	GameInfo      service.GameStateResponse
	IsNewGame     bool
	IsPrivateGame bool
	Loading       bool
}

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

	state = &State{
		CurrentStep: &StepNode{
			BreadCrumb: "Start",
			Previous:   nil,
		},
	}

	mainModel = *NewMainModel(state)
	askKeyModel = *NewAskKeyModel(state)
	askNameModel = *NewAskNameModel(state)
	createGameModel = *NewCreateGameModel(state)
	exitModel = *NewExitModel(state)
	playGameModel = *NewPlayGameModel(state)
	selectGameModel = *NewSelectGameModel(state)
	startOrJoinModel = *NewStartOrJoinModel(state)

	state.CurrentModel = mainModel

	return &mainModel
}

func (s *State) CantContinueModel(message string) (tea.Model, tea.Cmd) {
	exitModel.message = message
	return exitModel, tea.Quit
}

func (s *State) SkipAskName() {
	log.Printf("Skip name step, argument passed: -name %s\n", s.PlayerName)
	s.CurrentModel = startOrJoinModel
}

func (s *State) SkipAskKey() {
	log.Printf("Skip model mode/ model select steps, argument passed: -key %s\n", s.Key)
	s.CurrentModel = playGameModel
}

func (s *State) PreviousModel() (tea.Model, tea.Cmd) {
	var prevModel ConnectFourModel = askNameModel
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
	s.NavigateBackward(prevModel)
	return prevModel, prevCmd
}

func (s *State) NextModel() (tea.Model, tea.Cmd) {

	var nextModel ConnectFourModel = mainModel
	var nextCmd tea.Cmd

	switch (s.CurrentModel).(type) {
	case MainModel:
		nextModel = askNameModel
	case AskNameModel:
		nextModel = startOrJoinModel
	case AskKeyModel:
		nextModel = playGameModel
		nextCmd = joinGame(s.Key, s.PlayerName)
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
		nextCmd = LoadGameInfo(s.Key)
	case SelectGameModel:
		log.Printf("Player selected model %s, starting model...\n", s.Key)
		nextModel = playGameModel
		nextCmd = joinGame(s.Key, s.PlayerName)
	}

	log.Printf("[Next] Current Model = %T Next Model = %T\n", s.CurrentModel, nextModel)
	s.NavigateForward(nextModel)
	return nextModel, nextCmd
}

func (s *State) NavigateBackward(to ConnectFourModel) {
	s.CurrentStep = s.CurrentStep.Previous
	s.CurrentStep.Model = &to
	s.CurrentStep.BreadCrumb = to.BreadCrumb()
	s.CurrentModel = to
}

func (s *State) NavigateForward(to ConnectFourModel) {
	s.CurrentModel = to
	s.CurrentStep = &StepNode{
		Previous:   s.CurrentStep,
		BreadCrumb: s.CurrentModel.BreadCrumb(),
	}
	to.Init()
}

// CommonUpdate handles the common update messages and key presses (enter, esc) - it returns Model and Cmd in
// the same way that Update returns plus a boolean to indicate the calling Update can immediately return.
func (s *State) CommonUpdate(_ tea.Msg) (tea.Model, tea.Cmd, bool) {
	return nil, nil, false
}

func (s *State) CommonView(detail string) string {
	return lipgloss.JoinVertical(lipgloss.Left,
		s.renderHeader(),
		styles.Page.Render(detail),
	)
}

func (s *State) renderHeader() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		styles.AppTitle.Render("ConnectFour 0.2"),
		s.renderBreadCrumb(),
	)
}

func (s *State) renderBreadCrumb() string {
	var out string

	for n := s.CurrentStep; n.Previous != nil; n = n.Previous {
		out = styles.BreadCrumbSeparator.Render(" > ") + styles.BreadCrumb.Render(n.BreadCrumb) + out
	}
	return out
}

type GameInfoMsg struct {
	info         service.GameStateResponse
	errorMessage string
}

// LoadGameInfo returns a Cmd that sends a GameInfoMsg when Game Info is fetched.
func LoadGameInfo(key string) tea.Cmd {
	return func() tea.Msg {
		info, err := backend.GameInfo(key)
		msg := GameInfoMsg{
			info: info,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This model key could not be found"
		} else {
			msg.info = info
		}
		return msg
	}
}

// joinGame returns a Cmd that sends a tea.Msg when model is joined.
func joinGame(key string, name string) tea.Cmd {
	return func() tea.Msg {
		info, err := backend.Join(key, name)
		msg := GameInfoMsg{
			info: info,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This model key could not be found"
		} else {
			msg.info = info
		}
		return msg
	}
}
