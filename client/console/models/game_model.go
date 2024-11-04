package models

import (
	"connectfour/game"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
	board         game.Board
	currentPlayer game.Disc
	selectedCol   int
	redColor      lipgloss.Style
	yellowColor   lipgloss.Style
}

func NewGameModel() GameModel {
	m := GameModel{}
	m.currentPlayer = game.RedDisc
	m.redColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	m.yellowColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))

	return m
}

func (m GameModel) Init() tea.Cmd {
	return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Reset the board
		case "r":
			m.board.Reset()
			m.currentPlayer = game.RedDisc

		// control which column to drop in
		case "left", "j":
			if m.selectedCol > 0 {
				m.selectedCol--
			}

		case "right", "l":
			if m.selectedCol < game.BoardWidth-1 {
				m.selectedCol++
			}

		// The "enter" key and the space-bar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			if m.board.AddDisc(m.selectedCol, m.currentPlayer) {
				m.currentPlayer = m.getNextPlayer()
				m.selectedCol = 0
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m GameModel) getNextPlayer() game.Disc {
	if m.currentPlayer == game.RedDisc {
		return game.YellowDisc
	}
	return game.RedDisc
}

func (m GameModel) View() string {

	grey := lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	b := strings.Builder{}

	b.WriteString(strings.Repeat(" ", (m.selectedCol*2)+1))
	b.WriteString(m.renderDiscWithColor(m.currentPlayer))
	b.WriteByte('\n')
	for row := 0; row < game.BoardHeight; row++ {
		b.WriteString(grey.Render("|"))
		for col := 0; col < game.BoardWidth; col++ {
			b.WriteString(m.renderDiscWithColor(m.board.Cell(row, col)))
			b.WriteString(grey.Render("|"))
		}
		b.WriteString("\n")
	}

	if m.board.HasConnectFour() {
		b.WriteString("Connect four!\n")
	}

	return b.String()
}

func (m GameModel) renderDiscWithColor(disc game.Disc) string {
	s := m.yellowColor
	if disc == game.RedDisc {
		s = m.redColor
	}
	return s.Render(string(disc.Render()))
}
