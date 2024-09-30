package game

type Disc int

const (
	NoDisc = iota
	RedDisc
	YellowDisc
)

func (d Disc) Render() byte {
	switch d {
	case NoDisc:
		return ' '
	case RedDisc:
		return 'X'
	case YellowDisc:
		return 'O'
	}
	return '/'
}

func NewDisc(input byte) Disc {
	switch input {
	case 'X':
		return RedDisc
	case 'O':
		return YellowDisc
	default:
		return NoDisc
	}
}
