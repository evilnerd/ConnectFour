package models

import (
	"connectfour/server"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AskKeySubStatus string

const (
	WaitingForInput   AskKeySubStatus = "waiting_for_input"
	ValidatingGameKey AskKeySubStatus = "validating_game_key"
	ShowError         AskKeySubStatus = "show_error"
)

type AskKeyModel struct {
	MainModel
	InputGameKey string
	ErrorMessage string
	SubStatus    AskKeySubStatus
	Text         textinput.Model // the input box
}

func NewAskKeyModel() AskKeyModel {
	return AskKeyModel{
		Text:      textinput.New(),
		SubStatus: WaitingForInput,
	}
}

func (m AskKeyModel) Init() tea.Cmd {
	return nil
}

func (m AskKeyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEscape:
			break
		case tea.KeyEnter:
			break
		}
	}
}

func (m AskKeyModel) View() string {
	view := ""
	if m.InputGameKey == "" && m.ErrorMessage == "" {
		view = styles.Label.Render("Enter the private game key that was shared with you")
		view += "\n" + m.Text.View()
	} else {
		view += styles.Label.Render("Game key ") + styles.Value.Render(m.InputGameKey) + "\n"
		if m.ErrorMessage != "" {
			view += styles.Label.Render(m.ErrorMessage) + "\n"
		} else if m.GameStatus == server.Unknown {
			view += styles.Subdued.Render("Looking for that game...\n")
		} else {
			view += styles.Subdued.Render("The game status is " + string(m.GameStatus) + "\n")
		}
	}

	return view
}
