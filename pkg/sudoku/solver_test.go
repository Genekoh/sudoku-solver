package sudoku_test

import (
	"testing"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

// solverTestCase holds a named puzzle along with whether it is expected to be
// solvable. Puzzles are stored as raw [81]int so each solver run can build a
// fresh *Board (solvers mutate the board in place).
type solverTestCase struct {
	name     string
	puzzle   [81]int
	solvable bool
}

var solverTestCases = []solverTestCase{
	{
		name: "easy",
		puzzle: [81]int{
			0, 0, 0, 2, 6, 0, 7, 0, 1,
			6, 8, 0, 0, 7, 0, 0, 9, 0,
			1, 9, 0, 0, 0, 4, 5, 0, 0,
			8, 2, 0, 1, 0, 0, 0, 4, 0,
			0, 0, 4, 6, 0, 2, 9, 0, 0,
			0, 5, 0, 0, 0, 3, 0, 2, 8,
			0, 0, 9, 3, 0, 0, 0, 7, 4,
			0, 4, 0, 0, 5, 0, 0, 3, 6,
			7, 0, 3, 0, 1, 8, 0, 0, 0,
		},
		solvable: true,
	},
	{
		name: "medium",
		puzzle: [81]int{
			0, 2, 0, 6, 0, 8, 0, 0, 0,
			5, 8, 0, 0, 0, 9, 7, 0, 0,
			0, 0, 0, 0, 4, 0, 0, 0, 0,
			3, 7, 0, 0, 0, 0, 5, 0, 0,
			6, 0, 0, 0, 0, 0, 0, 0, 4,
			0, 0, 8, 0, 0, 0, 0, 1, 3,
			0, 0, 0, 0, 2, 0, 0, 0, 0,
			0, 0, 9, 8, 0, 0, 0, 3, 6,
			0, 0, 0, 3, 0, 6, 0, 9, 0,
		},
		solvable: true,
	},
	{
		name: "hard",
		puzzle: [81]int{
			0, 0, 0, 6, 0, 0, 4, 0, 0,
			7, 0, 0, 0, 0, 3, 6, 0, 0,
			0, 0, 0, 0, 9, 1, 0, 8, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 5, 0, 1, 8, 0, 0, 0, 3,
			0, 0, 0, 3, 0, 6, 0, 4, 5,
			0, 4, 0, 2, 0, 0, 0, 6, 0,
			9, 0, 3, 0, 0, 0, 0, 0, 0,
			0, 2, 0, 0, 0, 0, 1, 0, 0,
		},
		solvable: true,
	},
	{
		name: "already-solved",
		puzzle: [81]int{
			5, 3, 4, 6, 7, 8, 9, 1, 2,
			6, 7, 2, 1, 9, 5, 3, 4, 8,
			1, 9, 8, 3, 4, 2, 5, 6, 7,
			8, 5, 9, 7, 6, 1, 4, 2, 3,
			4, 2, 6, 8, 5, 3, 7, 9, 1,
			7, 1, 3, 9, 2, 4, 8, 5, 6,
			9, 6, 1, 5, 3, 7, 2, 8, 4,
			2, 8, 7, 4, 1, 9, 6, 3, 5,
			3, 4, 5, 2, 8, 6, 1, 7, 9,
		},
		solvable: true,
	},
	{
		name: "nearly-empty",
		puzzle: [81]int{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1,
		},
		solvable: true,
	},
	{
		name: "unsolvable",
		// Two 1s in the top row makes this puzzle impossible to complete.
		puzzle: [81]int{
			1, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		solvable: false,
	},
}

// runSolverTests exercises a solver against the shared set of test cases. Any
// implementation of sudoku.SudokuSolver can be validated by passing it here.
func runSolverTests(t *testing.T, solve sudoku.SudokuSolver) {
	t.Helper()

	for _, tc := range solverTestCases {
		t.Run(tc.name, func(t *testing.T) {
			b := sudoku.NewBoard(tc.puzzle)
			solved := solve(b)

			if tc.solvable {
				if !solved {
					t.Fatalf("expected puzzle to be solvable, but solver reported no solution")
				}
				if !b.IsCompleted() {
					t.Fatalf("solver reported success but board is not a completed valid solution")
				}
			} else {
				if solved {
					t.Fatalf("expected puzzle to be unsolvable, but solver reported a solution")
				}
			}
		})
	}
}

func TestNaiveBacktrackSolve(t *testing.T) {
	runSolverTests(t, sudoku.NaiveBacktrackSolve)
}

func TestBacktrackSolve(t *testing.T) {
	runSolverTests(t, sudoku.BacktrackSolve)
}
