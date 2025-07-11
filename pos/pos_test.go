package pos

import "testing"

const (
	_ uint8 = iota
	moveLine
	moveColumn
)

func Test_Move(t *testing.T) {
	tests := []struct {
		cmds     []uint8
		expected *pos
	}{
		{expected: newPos(1, 1)},
		{
			cmds:     []uint8{moveColumn},
			expected: newPos(1, 2),
		},
		{
			cmds:     []uint8{moveColumn, moveColumn, moveColumn},
			expected: newPos(1, 4),
		},
		{
			cmds:     []uint8{moveLine},
			expected: newPos(2, 1),
		},
		{
			cmds:     []uint8{moveLine, moveLine, moveLine},
			expected: newPos(4, 1),
		},
		{
			cmds: []uint8{
				moveColumn,
				moveColumn,
				moveLine,
				moveColumn,
				moveLine,
				moveLine,
				moveColumn,
				moveColumn,
				moveColumn,
			},
			expected: newPos(4, 4),
		},
	}

	for i, test := range tests {
		p := newPos(1, 1)

		for _, cmd := range test.cmds {
			if cmd == moveColumn {
				p.MoveColumn()
				continue
			}
			p.MoveLine()
		}

		if p.String() != test.expected.String() {
			t.Errorf("%d: ожидалось %s, получено %s", i, test.expected, p)
		}
	}
}
