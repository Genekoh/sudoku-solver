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
			b, err := sudoku.NewBoard(tc.puzzle)
			if err != nil {
				t.Fatalf("build test board: %v", err)
			}
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

func TestBacktrackMRVSolve(t *testing.T) {
	runSolverTests(t, sudoku.BacktrackMRVSolve)
}

func TestDancingLinksSolve(t *testing.T) {
	runSolverTests(t, sudoku.DancingLinksSolve)
}

func TestSolversRejectFilledInvalidBoardWithoutMutation(t *testing.T) {
	invalid := [81]int{
		1, 1, 3, 4, 5, 6, 7, 8, 9,
		4, 5, 6, 7, 8, 9, 1, 2, 3,
		7, 8, 9, 1, 2, 3, 4, 5, 6,
		2, 3, 4, 5, 6, 7, 8, 9, 1,
		5, 6, 7, 8, 9, 1, 2, 3, 4,
		8, 9, 1, 2, 3, 4, 5, 6, 7,
		3, 4, 5, 6, 7, 8, 9, 1, 2,
		6, 7, 8, 9, 1, 2, 3, 4, 5,
		9, 2, 5, 3, 4, 8, 6, 7, 1,
	}

	for name, solve := range map[string]sudoku.SudokuSolver{
		"naive":     sudoku.NaiveBacktrackSolve,
		"backtrack": sudoku.BacktrackSolve,
		"mrv":       sudoku.BacktrackMRVSolve,
	} {
		t.Run(name, func(t *testing.T) {
			board, err := sudoku.NewBoard(invalid)
			if err != nil {
				t.Fatalf("NewBoard: %v", err)
			}
			before := *board

			if solve(board) {
				t.Fatal("solver accepted a filled invalid board")
			}
			if *board != before {
				t.Fatal("solver mutated board after failure")
			}
		})
	}
}

// runSolverBenchmarks measures a solver against each puzzle in the shared set.
// A fresh board is built per iteration (solvers mutate in place); the timer is
// stopped during setup so only solve time is measured.
func runSolverBenchmarks(b *testing.B, solve sudoku.SudokuSolver) {
	for _, tc := range solverTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				board, err := sudoku.NewBoard(tc.puzzle)
				if err != nil {
					b.Fatalf("build benchmark board: %v", err)
				}
				b.StartTimer()

				solve(board)
			}
		})
	}
}

func BenchmarkNaiveBacktrackSolve(b *testing.B) {
	runSolverBenchmarks(b, sudoku.NaiveBacktrackSolve)
}

func BenchmarkBacktrackSolve(b *testing.B) {
	runSolverBenchmarks(b, sudoku.BacktrackSolve)
}

func BenchmarkBacktrackMRVSolve(b *testing.B) {
	runSolverBenchmarks(b, sudoku.BacktrackMRVSolve)
}

func BenchmarkDancingLinksSolve(b *testing.B) {
	runSolverBenchmarks(b, sudoku.DancingLinksSolve)
}
