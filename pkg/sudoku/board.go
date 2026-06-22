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

func NewBoard(nums [81]int) *Board {
	var b Board
	for i, n := range nums {
		b[i/9][i%9] = Digit(n)
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
			if cell != Empty {
				rowStr += fmt.Sprintf("%d", cell)
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
