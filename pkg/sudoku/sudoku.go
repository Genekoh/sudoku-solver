package sudoku

import (
	"slices"
)

func isValidSudokuSet(set [9]Cell) bool {
	seen := []int{}
	for _, c := range set {
		if c.Number == Empty {
			continue
		}

		if slices.Contains(seen, c.Number) {
			return false
		} else {
			seen = append(seen, c.Number)
		}
	}
	return true
}

func (b Board) IsValid() bool {
	// Validate Each Row
	for _, row := range b {
		if !isValidSudokuSet(row) {
			return false
		}
	}

	// Validate Each Column
	for i := 0; i < 9; i++ {
		col := [9]Cell{}
		for j, row := range b {
			col[j] = row[i]
		}

		if !isValidSudokuSet(col) {
			return false
		}
	}

	// Validate Each Block
	for m := 0; m < 3; m++ {
		for n := 0; n < 3; n++ {
			block := [9]Cell{}
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					block[3*i+j] = b[3*m+i][3*n+j]
				}
			}

			if !isValidSudokuSet(block) {
				return false
			}
		}
	}

	return true
}

func (b Board) IsFilled() bool {
	for _, row := range b {
		for _, c := range row {
			if c.Number == Empty {
				return false
			}
		}
	}

	return true
}

func (b Board) IsCompleted() bool {
	return b.IsValid() && b.IsFilled()
}

type SudokuSolver func(*Board) bool

func BacktrackSolve(b *Board) bool {
	if b.IsCompleted() {
		return true
	}

	// Traverse to the first empty cell found
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if b[i][j].Number != Empty {
				continue
			}

			// For each possible
			for n := 1; n <= 9; n++ {
				b[i][j].Number = n
				if b.IsValid() {
					solved := BacktrackSolve(b)
					if solved {
						return true
					}
				}
				b[i][j].Number = Empty
			}

			// No Solutions Found
			return false
		}
	}

	panic("Invalid Sudoku Given")
}
