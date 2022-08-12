package main

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
	assert.False(t, b.hasConnectFour())
}

func TestBoard_hasConnectFour_FindsHorizontal(t *testing.T) {

	// 1. bottom left
	b := getTestBoard()
	b.addDisc(0, RedDisc)
	b.addDisc(1, RedDisc)
	b.addDisc(2, RedDisc)
	b.addDisc(3, RedDisc)

	// Act & Assert
	assert.True(t, b.hasConnectFour())

	// 2. bottom right
	b.reset()
	b.addDisc(BoardWidth-1, RedDisc)
	b.addDisc(BoardWidth-2, RedDisc)
	b.addDisc(BoardWidth-3, RedDisc)
	b.addDisc(BoardWidth-4, RedDisc)

	// Act & Assert
	assert.True(t, b.hasConnectFour())

}

func TestBoard_hasConnectFour_FindsVertical(t *testing.T) {

	b := getTestBoard()
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(0, RedDisc)

	// Act & Assert
	assert.True(t, b.hasConnectFour())

}

func TestBoard_HasConnectFour_FindsDiagonalForward(t *testing.T) {

	// Arrange
	b := getTestBoard()
	// | R | Y | R | Y |   |   |
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(2, RedDisc)
	b.addDisc(3, YellowDisc)

	// |   | R | Y | R |   |   |
	// | R | Y | R | Y | Y |   |
	b.addDisc(1, RedDisc)
	b.addDisc(2, YellowDisc)
	b.addDisc(3, RedDisc)
	b.addDisc(4, YellowDisc)

	// |   |   | R | Y |   |   |
	// |   | R | Y | R | R |   |
	// | R | Y | R | Y | Y | Y |
	b.addDisc(2, RedDisc)
	b.addDisc(3, YellowDisc)
	b.addDisc(4, RedDisc)
	b.addDisc(5, YellowDisc)

	// |   |   |   | R |   |   |
	// |   |   | R | Y |   |   |
	// |   | R | Y | R | R |   |
	// | R | Y | R | Y | Y | Y |
	b.addDisc(3, RedDisc)

	// Act & Assert
	assert.True(t, b.hasConnectFour())
}

func TestBoard_HasConnectFour_FindsDiagonalBackward(t *testing.T) {

	// Arrange
	b := getTestBoard()
	// | R | Y | R | Y |   |   |
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(2, RedDisc)
	b.addDisc(3, YellowDisc)

	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.addDisc(0, RedDisc)
	b.addDisc(1, RedDisc)
	b.addDisc(2, YellowDisc)
	b.addDisc(3, YellowDisc)

	// | R | Y | R |   |   |   |
	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.addDisc(0, RedDisc)
	b.addDisc(1, YellowDisc)
	b.addDisc(2, RedDisc)

	// | Y |   |   |   |   |   |
	// | R | Y | R |   |   |   |
	// | R | R | Y | Y |   |   |
	// | R | Y | R | Y |   |   |
	b.addDisc(0, YellowDisc)

	// Act & Assert
	log.Println("\n" + b.String())
	assert.True(t, b.hasConnectFour())
}
