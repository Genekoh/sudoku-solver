package sudoku

import (
	"fmt"
	"strings"
)

type Digit uint8

const (
	Empty Digit = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

type Board [9][9]Digit

// NewBoard constructs a board from row-major values. Every value must be in
// the range 0 through 9, where 0 represents an empty cell.
func NewBoard(nums [81]int) (*Board, error) {
	var b Board
	for i, n := range nums {
		if n < int(Empty) || n > int(Nine) {
			return nil, fmt.Errorf("cell %d: digit %d is outside the range 0..9", i, n)
		}
		b[i/9][i%9] = Digit(n)
	}

	return &b, nil
}

// String returns a human-readable rendering of the board. Empty cells are
// represented by x.
func (b Board) String() string {
	var output strings.Builder
	for i, row := range b {
		var rowStr strings.Builder
		for j, cell := range row {
			if j == 3 || j == 6 {
				rowStr.WriteString(" |")
			}

			if j != 0 {
				rowStr.WriteByte(' ')
			}
			if cell != Empty {
				fmt.Fprint(&rowStr, cell)
			} else {
				rowStr.WriteByte('x')
			}
		}

		if i == 3 || i == 6 {
			output.WriteString(strings.Repeat("-", rowStr.Len()))
			output.WriteByte('\n')
		}
		output.WriteString(rowStr.String())
		if i != len(b)-1 {
			output.WriteByte('\n')
		}
	}
	return output.String()
}
