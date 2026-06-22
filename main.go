package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

// solver bundles a human-readable name with the function that solves a board
// in place, returning true on success.
type solver struct {
	name string
	fn   func(*sudoku.Board) bool
}

var solvers = map[string]solver{
	"naive":     {"Naive backtracking", sudoku.NaiveBacktrackSolve},
	"backtrack": {"Backtracking", sudoku.BacktrackSolve},
	"mrv":       {"MRV backtracking", sudoku.BacktrackMRVSolve},
}

// defaultBoard is a sample puzzle used when no board is provided.
const defaultBoard = "210903000000071800000002006007000013600000057001000000062050000000600400304000000"

func main() {
	solverFlag := flag.String("solver", "", "solver to use: naive, backtrack, mrv")
	boardFlag := flag.String("board", "", "starting board as 81 chars (0 or . for empty); use '-' to read from stdin")
	flag.Parse()

	solverKey := *solverFlag
	if solverKey == "" {
		solverKey = chooseSolver()
	}
	s, ok := solvers[solverKey]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown solver %q (choose from: naive, backtrack, mrv)\n", solverKey)
		os.Exit(1)
	}

	raw := *boardFlag
	if raw == "" || raw == "-" {
		raw = chooseBoard()
	}

	board, err := parseBoard(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid board: %v\n", err)
		os.Exit(1)
	}

	if !board.IsValid() {
		fmt.Fprintln(os.Stderr, "invalid board: starting position breaks sudoku rules")
		os.Exit(1)
	}

	fmt.Printf("\nSolver: %s\n\nGiven:\n", s.name)
	sudoku.PrintBoard(*board)

	start := time.Now()
	solved := s.fn(board)
	elapsed := time.Since(start)

	fmt.Println()
	if !solved {
		fmt.Printf("Failed to solve (took %s)\n", elapsed)
		os.Exit(1)
	}
	fmt.Printf("Solved in %s:\n", elapsed)
	sudoku.PrintBoard(*board)
}

// chooseSolver interactively prompts for a solver key.
func chooseSolver() string {
	fmt.Println("Choose a solver:")
	fmt.Println("  1) naive     - Naive backtracking")
	fmt.Println("  2) backtrack - Backtracking")
	fmt.Println("  3) mrv       - MRV backtracking")
	fmt.Print("> ")

	switch strings.TrimSpace(readLine()) {
	case "1", "naive":
		return "naive"
	case "2", "backtrack":
		return "backtrack"
	case "3", "mrv":
		return "mrv"
	default:
		return ""
	}
}

// chooseBoard interactively prompts for a starting board, defaulting to the
// built-in sample when left blank.
func chooseBoard() string {
	fmt.Println("\nEnter a starting board as 81 chars (0 or . for empty),")
	fmt.Println("or press Enter to use the built-in sample puzzle.")
	fmt.Print("> ")

	raw := strings.TrimSpace(readLine())
	if raw == "" {
		return defaultBoard
	}
	return raw
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

// parseBoard converts an 81-character board string into a *Board. Both '0' and
// '.' denote empty cells; whitespace is ignored.
func parseBoard(raw string) (*sudoku.Board, error) {
	var cleaned strings.Builder
	for _, r := range raw {
		switch {
		case r >= '1' && r <= '9':
			cleaned.WriteRune(r)
		case r == '0' || r == '.':
			cleaned.WriteByte('0')
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			// ignore whitespace
		default:
			return nil, fmt.Errorf("unexpected character %q", r)
		}
	}

	s := cleaned.String()
	if len(s) != 81 {
		return nil, fmt.Errorf("expected 81 cells, got %d", len(s))
	}

	var nums [81]int
	for i, c := range s {
		nums[i] = int(c - '0')
	}
	return sudoku.NewBoard(nums), nil
}
