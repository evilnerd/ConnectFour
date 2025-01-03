package models

import (
	"connectfour/server"
	tea "github.com/charmbracelet/bubbletea"
)

// BackStepMsg is sent when the user pressed ESC.
type BackStepMsg struct{}

func BackCmd() tea.Msg {
	return BackStepMsg{}
}

// QuitMsg is sent when the user pressed ctrl-c
type QuitMsg struct{}

func QuitCmd() tea.Msg {
	return QuitMsg{}
}

type NotConnected struct {
	message string
}

type GameCreated struct {
	game server.NewGameResponse
}
type Connected struct{}
