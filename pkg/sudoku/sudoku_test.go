package sudoku_test

import (
	"testing"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

func TestBacktrackSolve(t *testing.T) {
	var solve_tests = [](*sudoku.Board){
		sudoku.NewBoard([81]int{
			0, 0, 0, 2, 6, 0, 7, 0, 1,
			6, 8, 0, 0, 7, 0, 0, 9, 0,
			1, 9, 0, 0, 0, 4, 5, 0, 0,
			8, 2, 0, 1, 0, 0, 0, 4, 0,
			0, 0, 4, 6, 0, 2, 9, 0, 0,
			0, 5, 0, 0, 0, 3, 0, 2, 8,
			0, 0, 9, 3, 0, 0, 0, 7, 4,
			0, 4, 0, 0, 5, 0, 0, 3, 6,
			7, 0, 3, 0, 1, 8, 0, 0, 0,
		}),
	}

	for _, b := range solve_tests {
		solved := sudoku.BacktrackSolve(b)
		if !solved || !b.IsCompleted() {
			t.Error("Should be able to solve")
		}
	}
}
