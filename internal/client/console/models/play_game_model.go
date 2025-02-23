package models

import (
	"connectfour/internal/client/console/backend"
	game2 "connectfour/internal/model"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
	"time"
)

type PlayGameModel struct {
	*State
	Loading       bool
	board         game2.Board
	currentPlayer game2.Disc
	selectedCol   int
	redColor      lipgloss.Style
	yellowColor   lipgloss.Style
}

type RefreshTickMsg time.Time

func doTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return RefreshTickMsg(t)
	})
}

func NewPlayGameModel(state *State) *PlayGameModel {
	m := &PlayGameModel{
		State: state,
	}
	m.currentPlayer = game2.RedDisc
	m.redColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	m.yellowColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	m.Loading = true
	return m
}

func (m PlayGameModel) BreadCrumb() string {
	return "Play"
}

func (m PlayGameModel) PlayMoveCmd(column int) tea.Cmd {
	m.Loading = true
	return func() tea.Msg {
		info, err := backend.Move(m.wc, m.Key, column)
		msg := GameInfoMsg{
			info: info,
		}
		if errors.Is(err, backend.GameNotFoundError{}) {
			msg.errorMessage = "This game key could not be found"
		}
		return msg
	}
}

// playing indicates whether we are really playing (true) or whether we are waiting for an event, or perhaps the model has ended.
func (m PlayGameModel) playing() bool {
	return m.Loading == false &&
		m.GameInfo.Status == game2.Started
}

func (m PlayGameModel) myTurn() bool {
	return m.GameInfo.PlayerTurnName == m.PlayerName
}

func (m PlayGameModel) Init() tea.Cmd {
	log.Printf("Init for PlayGameModel - getting game data and starting ticker\n")
	return tea.Batch(LoadGameInfo(m.Key), doTick())
}

func (m PlayGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case RefreshTickMsg:
		//log.Printf("Refresh tick %s\n", time.Time(msg).Format("01-02-2006 15:04:05.000000"))
		return m, tea.Batch(LoadGameInfo(m.Key), doTick())

	case GameInfoMsg:
		m.GameInfo = msg.info
		m.board = *game2.FromMap(m.GameInfo.Board)
		m.currentPlayer = game2.Disc(m.GameInfo.PlayerTurn)
		m.Loading = false
		if msg.errorMessage != "" {
			log.Printf("There was an error getting the game state: %s\n", msg.errorMessage)
			return m.PreviousModel()
		}

	// Is it a key press?
	case tea.KeyMsg:
		if m.playing() {
			switch msg.String() {
			case "esc", "ctrl+c", "q":
				return m.PreviousModel()

			// control which column to drop in
			case "left", "j":
				if m.selectedCol > 0 {
					m.selectedCol--
				}

			case "right", "l":
				if m.selectedCol < game2.BoardWidth-1 {
					m.selectedCol++
				}

			case "enter", " ":
				if m.myTurn() {
					return m, m.PlayMoveCmd(m.selectedCol + 1) // the api expects 1-7 for columns.
				} else {
					return m, LoadGameInfo(m.Key)
				}
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

	view := ""
	if m.GameInfo.Status == "" {
		view = styles.Header.Render("Updating game info...")
	} else if m.GameInfo.Status == game2.Started {
		view = m.renderGameBoard()
	} else if m.GameInfo.Status == game2.Created {
		view = styles.Header.Render("Waiting for other player... come back later.")
	} else if m.GameInfo.Status == game2.Finished {
		view = styles.Header.Render("This game has finished. Better create a new one.")
	} else {
		view = styles.Header.Render("The game is no longer valid. ")
	}

	return m.CommonView(view)

}

func (m PlayGameModel) renderGameBoard() string {
	grey := lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	b := strings.Builder{}

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		// Playing as
		styles.Label.Render("Playing a game as ")+styles.Value.Render(m.PlayerName),
		// Key and players
		styles.Subdued.Render("Gamekey ")+
			styles.Value.Render(m.Key)+
			styles.Subdued.Render(", game between ")+
			styles.Value.Render(m.GameInfo.Player1Name)+
			styles.Subdued.Render(" and ")+
			styles.Value.Render(m.GameInfo.Player2Name),
		// Whose turn is it
		styles.Label.Render("Player turn: ")+styles.Value.Render(m.GameInfo.PlayerTurnName),
	))

	b.WriteRune('\n')
	b.WriteRune('\n')
	if m.myTurn() {
		b.WriteString(strings.Repeat(" ", (m.selectedCol*2)+1))
		b.WriteString(m.renderDiscWithColor(m.currentPlayer))
		b.WriteRune('\n')
	} else {
		b.WriteString(styles.Label.Render("Waiting for other player move"))
		b.WriteRune('\n')
	}
	for row := 0; row < game2.BoardHeight; row++ {
		b.WriteString(grey.Render("|"))
		for col := 0; col < game2.BoardWidth; col++ {
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

func (m PlayGameModel) renderDiscWithColor(disc game2.Disc) string {
	s := m.yellowColor
	if disc == game2.RedDisc {
		s = m.redColor
	}
	return s.Render(string(disc.Render()))
}
