package models

import (
	"connectfour/client/console/backend"
	"connectfour/game"
	"connectfour/server"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
)

type PlayGameModel struct {
	*State
	board         game.Board
	currentPlayer game.Disc
	selectedCol   int
	redColor      lipgloss.Style
	yellowColor   lipgloss.Style
}

func NewPlayGameModel(state *State) *PlayGameModel {
	m := &PlayGameModel{
		State: state,
	}
	m.currentPlayer = game.RedDisc
	m.redColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	m.yellowColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))

	return m
}

type PlayMsg struct {
	info         server.GameStateResponse
	errorMessage string
}

func (m PlayGameModel) PlayMoveCmd(column int) tea.Cmd {
	m.Loading = true
	return func() tea.Msg {
		info, err := backend.Move(m.Key, column)
		msg := PlayMsg{
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

// playing indicates whether we are really playing (true) or whether we are waiting for an event, or perhaps the game has ended.
func (m PlayGameModel) playing() bool {
	return m.Loading == false &&
		m.GameInfo.Status == server.Started
}

func (m PlayGameModel) Init() tea.Cmd {
	return loadGameInfo(m.Key)
}

func (m PlayGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case GameInfoMsg:
		m.GameInfo = msg.info
		if msg.errorMessage != "" {
			log.Printf("There was an error getting the game state: %s\n", msg.errorMessage)
			return m.PreviousModel()
		}

	case PlayMsg:
		m.GameInfo = msg.info
		m.board = *game.FromMap(m.GameInfo.Board)
		m.currentPlayer = game.Disc(m.GameInfo.PlayerTurn)
		m.Loading = false
		if msg.errorMessage != "" {
			log.Printf("There was an error returned when trying to play the move: %s\n", msg.errorMessage)
			return m.PreviousModel()
		}

	// Is it a key press?
	case tea.KeyMsg:
		if m.playing() {
			switch msg.String() {
			case "esc", "ctrl+c", "q":
				return m.PreviousModel()

			// Reset the board
			//case "r":
			//	m.board.Reset()
			//	m.currentPlayer = game.RedDisc

			// control which column to drop in
			case "left", "j":
				if m.selectedCol > 0 {
					m.selectedCol--
				}

			case "right", "l":
				if m.selectedCol < game.BoardWidth-1 {
					m.selectedCol++
				}

			case "enter", " ":
				return m, m.PlayMoveCmd(m.selectedCol + 1) // the server expects 1-7 for columns.
			}
		} else {
			if m.GameInfo.Status != "" {
				return m.PreviousModel()
			}
		}
	}

	return m, nil
}

func (m PlayGameModel) View() string {

	if m.GameInfo.Status == "" {
		return styles.Header.Render("Updating game info...")
	} else if m.GameInfo.Status == server.Started {
		grey := lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
		b := strings.Builder{}

		b.WriteString(styles.Header.Render("Playing a game\n"))
		b.WriteString(styles.Subdued.Render("Gamekey "))
		b.WriteString(styles.Value.Render(m.Key))
		b.WriteString(styles.Subdued.Render(", game between "))
		b.WriteString(styles.Value.Render(m.GameInfo.Player1Name))
		b.WriteString(styles.Subdued.Render(" and "))
		b.WriteString(styles.Value.Render(m.GameInfo.Player2Name))
		b.WriteRune('\n')
		b.WriteRune('\n')

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
	} else if m.GameInfo.Status == server.Created {
		return styles.Header.Render("Waiting for other player... come back later.")
	} else if m.GameInfo.Status == server.Finished {
		return styles.Header.Render("This game has finished. Better create a new one.")
	} else {
		return styles.Header.Render("The game is no longer valid. ")
	}
}

func (m PlayGameModel) renderDiscWithColor(disc game.Disc) string {
	s := m.yellowColor
	if disc == game.RedDisc {
		s = m.redColor
	}
	return s.Render(string(disc.Render()))
}
