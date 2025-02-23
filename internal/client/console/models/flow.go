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
	CurrentStep        *StepNode
	CurrentModel       ConnectFourModel
	Key                string
	PlayerName         string
	PlayerEmail        string
	GameStatus         model.GameStatus
	GameInfo           service.GameStateResponse
	IsNewGame          bool
	IsPrivateGame      bool
	IsContinue         bool // When the game mode is to continue a running game.
	MustReauthenticate bool // Set when the JWT expires or is invalid somehow.
	NoAuthStorage      bool // Set as a cmd arg flag to indicate we should not load nor save the JWT (for testing)
	wc                 *backend.WebClient
}

var (
	wc               *backend.WebClient
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

func CreateModels(key string, storeInFile bool) *MainModel {
	args := []backend.WebClientOption{
		backend.WithBaseUrl(backend.ServerUrl),
		backend.WithStoreInFile(storeInFile, backend.JwtFileName()),
		backend.WithReAuthCallback(func() {
			ReAuthenticate(state)
		}),
	}

	wc, _ = backend.NewWebClient(args...)
	state = &State{
		CurrentStep: &StepNode{
			BreadCrumb: "Start",
			Previous:   nil,
		},
		Key: key,
		wc:  wc,
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

func (s *State) SkipAskLogin() {
	log.Printf("Skip name step, argument passed: -name %s\n", s.PlayerName)
	s.CurrentModel = startOrJoinModel
}

func (s *State) SkipAskKey() {
	log.Printf("Skip game mode/ game select steps, argument passed: -key %s\n", s.Key)
	s.CurrentModel = playGameModel
}

func (s *State) PreviousModel() (tea.Model, tea.Cmd) {
	var prevModel ConnectFourModel = askNameModel
	var prevCmd tea.Cmd
	if s.MustReauthenticate {
		return prevModel, nil
	}
	switch (s.CurrentModel).(type) {
	case AskNameModel, StartOrJoinModel:
		prevModel = exitModel
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
	if s.MustReauthenticate {
		return askNameModel, nil
	}
	switch (s.CurrentModel).(type) {
	case MainModel:
		nextModel = askNameModel
	case AskNameModel:
		nextModel = startOrJoinModel
	case AskKeyModel:
		nextModel = playGameModel
		nextCmd = joinGame(s.Key)
	case StartOrJoinModel:
		if s.IsContinue {
			nextModel = selectGameModel
			nextCmd = selectGameModel.loadMyGames()
		} else if s.IsNewGame {
			nextModel = createGameModel
			nextCmd = createGame(s.IsPrivateGame)
		} else {
			if s.IsPrivateGame {
				nextModel = askKeyModel
			} else {
				nextModel = selectGameModel
				nextCmd = selectGameModel.loadOpenGames()
			}
		}
	case CreateGameModel:
		nextModel = playGameModel
		nextCmd = LoadGameInfo(s.Key)
	case SelectGameModel:
		log.Printf("Player selected game %s, starting game...\n", s.Key)
		nextModel = playGameModel
		nextCmd = joinGame(s.Key)
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

// ReAuthenticate sets a flag to indicate that the next model should always be the [AskNameModel]
func ReAuthenticate(state *State) {
	state.PlayerName = ""
	state.PlayerEmail = ""
	state.MustReauthenticate = true
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

	errmsg := ""

	if s.MustReauthenticate {
		errmsg = styles.Label.Render("The server says that you need to authenticate again. So I'm going to ask for your credentials now.\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.AppTitle.Render("ConnectFour 0.2"),
		styles.Error.Render(errmsg),
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
		info, err := backend.GameInfo(wc, key)
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
func joinGame(key string) tea.Cmd {
	return func() tea.Msg {
		info, err := backend.Join(wc, key)
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
