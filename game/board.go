package game

import (
	strings "strings"
)

const (
	BoardWidth  = 7
	BoardHeight = 6
)

type Board struct {
	cells [BoardWidth * BoardHeight]Disc
}

func FromMap(boardMap map[int]string) *Board {
	board := &Board{}
	for row, values := range boardMap {
		for col := 0; col < BoardWidth; col++ {
			board.setCell(row-1, col, NewDisc(values[col]))
		}
	}

	return board
}

func (b *Board) cellIndex(row int, col int) int {
	if row < 0 || row > BoardHeight-1 {
		panic("value of invalid row requested")
	}

	if col < 0 || col > BoardWidth-1 {
		panic("value of invalid column requested")
	}

	return row*BoardWidth + col
}
func (b *Board) Cell(row int, col int) Disc {
	return b.cells[b.cellIndex(row, col)]
}

func (b *Board) setCell(row int, col int, d Disc) {
	b.cells[b.cellIndex(row, col)] = d
}

func (b *Board) AddDisc(col int, disc Disc) bool {
	for row := BoardHeight - 1; row >= 0; row-- {
		if b.Cell(row, col) == NoDisc {
			b.setCell(row, col, disc)
			return true
		}
	}
	return false
}

func (b *Board) Reset() {
	b.cells = [BoardWidth * BoardHeight]Disc{}
}

func (b *Board) HasConnectFour() bool {
	// horizontal
	for row := 0; row < BoardHeight; row++ {
		for col := 0; col < BoardWidth-3; col++ {
			c := b.Cell(row, col)
			if c != NoDisc &&
				c == b.Cell(row, col+1) &&
				c == b.Cell(row, col+2) &&
				c == b.Cell(row, col+3) {
				return true
			}
		}
	}

	// vertical
	for col := 0; col < BoardWidth; col++ {
		for row := 0; row < BoardHeight-3; row++ {
			c := b.Cell(row, col)
			if c != NoDisc &&
				c == b.Cell(row+1, col) &&
				c == b.Cell(row+2, col) &&
				c == b.Cell(row+3, col) {
				return true
			}
		}
	}

	// Diagonal //
	for col := 3; col < BoardWidth; col++ {
		for row := 0; row < BoardHeight-3; row++ {
			c := b.Cell(row, col)
			if c != NoDisc &&
				c == b.Cell(row+1, col-1) &&
				c == b.Cell(row+2, col-2) &&
				c == b.Cell(row+3, col-3) {
				return true
			}
		}
	}

	// Diagonal \\
	for col := 0; col < BoardWidth-3; col++ {
		for row := 0; row < BoardHeight-3; row++ {
			c := b.Cell(row, col)
			if c != NoDisc &&
				c == b.Cell(row+1, col+1) &&
				c == b.Cell(row+2, col+2) &&
				c == b.Cell(row+3, col+3) {
				return true
			}
		}
	}

	return false
}

func (b *Board) String() string {
	sb := strings.Builder{}
	for row := 0; row < BoardHeight; row++ {
		sb.WriteString("|")
		for col := 0; col < BoardWidth; col++ {
			sb.WriteByte(b.Cell(row, col).Render())
			sb.WriteString("|")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *Board) Map() map[int]string {
	output := make(map[int]string)
	for row := 0; row < BoardHeight; row++ {
		sb := strings.Builder{}
		for col := 0; col < BoardWidth; col++ {
			sb.WriteByte(b.Cell(row, col).Render())
		}
		output[row+1] = sb.String()
	}
	return output
}
