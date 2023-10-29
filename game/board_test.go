package game

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func getTestBoard() Board {
	return Board{
		cells: [42]Disc{},
	}
}

func TestBoard_hasConnectFour_ReturnsFalseForEmptyBoard(t *testing.T) {
	// Arrange
	b := getTestBoard()

	// Act & Assert
	assert.False(t, b.HasConnectFour())
}

func TestBoard_hasConnectFour_FindsHorizontal(t *testing.T) {

	// 1. bottom left
	b := getTestBoard()
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, RedDisc)
	b.AddDisc(2, RedDisc)
	b.AddDisc(3, RedDisc)

	// Act & Assert
	assert.True(t, b.HasConnectFour())

	// 2. bottom right
	b.Reset()
	b.AddDisc(BoardWidth-1, RedDisc)
	b.AddDisc(BoardWidth-2, RedDisc)
	b.AddDisc(BoardWidth-3, RedDisc)
	b.AddDisc(BoardWidth-4, RedDisc)

	// Act & Assert
	assert.True(t, b.HasConnectFour())

}

func TestBoard_hasConnectFour_FindsVertical(t *testing.T) {

	b := getTestBoard()
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(0, RedDisc)

	// Act & Assert
	assert.True(t, b.HasConnectFour())

}

func TestBoard_HasConnectFour_FindsDiagonalForward(t *testing.T) {

	// Arrange
	b := getTestBoard()
	// | R | Y | R | Y |   |   |
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(2, RedDisc)
	b.AddDisc(3, YellowDisc)

	// |   | R | Y | R |   |   |
	// | R | Y | R | Y | Y |   |
	b.AddDisc(1, RedDisc)
	b.AddDisc(2, YellowDisc)
	b.AddDisc(3, RedDisc)
	b.AddDisc(4, YellowDisc)

	// |   |   | R | Y |   |   |
	// |   | R | Y | R | R |   |
	// | R | Y | R | Y | Y | Y |
	b.AddDisc(2, RedDisc)
	b.AddDisc(3, YellowDisc)
	b.AddDisc(4, RedDisc)
	b.AddDisc(5, YellowDisc)

	// |   |   |   | R |   |   |
	// |   |   | R | Y |   |   |
	// |   | R | Y | R | R |   |
	// | R | Y | R | Y | Y | Y |
	b.AddDisc(3, RedDisc)

	// Act & Assert
	assert.True(t, b.HasConnectFour())
}

func TestBoard_HasConnectFour_FindsDiagonalBackward(t *testing.T) {

	// Arrange
	b := getTestBoard()
	// | R | Y | R | Y |   |   |
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(2, RedDisc)
	b.AddDisc(3, YellowDisc)

	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, RedDisc)
	b.AddDisc(2, YellowDisc)
	b.AddDisc(3, YellowDisc)

	// | R | Y | R |   |   |   |
	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.AddDisc(0, RedDisc)
	b.AddDisc(1, YellowDisc)
	b.AddDisc(2, RedDisc)

	// | Y |   |   |   |   |   |
	// | R | Y | R |   |   |   |
	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.AddDisc(0, YellowDisc)

	// Act & Assert
	log.Println("\n" + b.String())
	assert.True(t, b.HasConnectFour())
}
