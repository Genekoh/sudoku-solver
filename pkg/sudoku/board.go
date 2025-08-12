package sudoku

import (
	"fmt"
	"strings"
)

const (
	Empty int = iota
	One       // consider just having Empty defined
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

type Cell struct {
	Number  int
	IsFixed bool
}

type Board [9][9]Cell

func NewBoard(nums [81]int) *Board {
	var b Board
	for i, n := range nums {
		isFixed := n != 0
		b[i/9][i%9] = Cell{n, isFixed}
	}

	return &b
}

func PrintBoard(b Board) {
	for i, row := range b {
		rowStr := ""
		for j, cell := range row {
			if j == 3 || j == 6 {
				rowStr += " |"
			}

			if j != 0 {
				rowStr += " "
			}
			if cell.Number != Empty {
				rowStr += fmt.Sprintf("%d", cell.Number)
			} else {
				rowStr += "x"
			}
		}

		if i == 3 || i == 6 {
			line := strings.Repeat("-", len(rowStr))
			fmt.Println(line)
		}
		fmt.Println(rowStr)
	}
}
