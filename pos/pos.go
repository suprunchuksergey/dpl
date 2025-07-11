package pos

import "fmt"

// отсчет начинается с 1
type pos struct {
	line   uint16
	column uint16
}

func newPos(line, column uint16) *pos {
	return &pos{
		line:   line,
		column: column,
	}
}

func (p *pos) MoveLine() {
	p.line++
	p.column = 1
}

func (p *pos) MoveColumn() {
	p.column++
}

func (p *pos) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (p *pos) Clone() Pos {
	return newPos(p.line, p.column)
}

func (p *pos) Line() uint16 {
	return p.line
}

func (p *pos) Column() uint16 {
	return p.column
}

type Pos interface {
	fmt.Stringer
	MoveLine()
	MoveColumn()
	Clone() Pos
	Line() uint16
	Column() uint16
}

func New() Pos {
	return NewWithStart(1, 1)
}

// line и column от 1
func NewWithStart(line, column uint16) Pos {
	return newPos(line, column)
}
