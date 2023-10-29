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
