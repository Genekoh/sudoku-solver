package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

// solver bundles a human-readable name with the function that solves a board
// in place, returning true on success.
type solver struct {
	key  string
	name string
	fn   func(*sudoku.Board) bool
}

var solvers = []solver{
	{"naive", "Naive backtracking", sudoku.NaiveBacktrackSolve},
	{"backtrack", "Backtracking", sudoku.BacktrackSolve},
	{"mrv", "MRV backtracking", sudoku.BacktrackMRVSolve},
	{"dlx", "Dancing links", sudoku.DancingLinksSolve},
}

// defaultBoard is a sample puzzle used when no board is provided.
const defaultBoard = "210903000000071800000002006007000013600000057001000000062050000000600400304000000"

func main() {
	solverFlag := flag.String("solver", "", fmt.Sprintf("solver to use: %s", solverChoices()))
	boardFlag := flag.String("board", "", "starting board as 81 chars (0 or . for empty); use '-' to read from stdin")
	flag.Parse()

	solverKey := *solverFlag
	if solverKey == "" {
		solverKey = chooseSolver()
	}
	s, ok := findSolver(solverKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown solver %q (choose from: %s)\n", solverKey, solverChoices())
		os.Exit(1)
	}

	raw, err := readBoardInput(*boardFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid board input: %v\n", err)
		os.Exit(1)
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
	fmt.Println(board)

	start := time.Now()
	solved := s.fn(board)
	elapsed := time.Since(start)

	fmt.Println()
	if !solved {
		fmt.Printf("Failed to solve (took %s)\n", elapsed)
		os.Exit(1)
	}
	fmt.Printf("Solved in %s:\n", elapsed)
	fmt.Println(board)
}

// chooseSolver interactively prompts for a solver key.
func chooseSolver() string {
	fmt.Println("Choose a solver:")
	for i, solver := range solvers {
		fmt.Printf("  %d) %-9s - %s\n", i+1, solver.key, solver.name)
	}
	fmt.Print("> ")

	choice := strings.TrimSpace(readLine())
	for i, solver := range solvers {
		if choice == fmt.Sprint(i+1) || choice == solver.key {
			return solver.key
		}
	}
	return ""
}

func findSolver(key string) (solver, bool) {
	for _, solver := range solvers {
		if solver.key == key {
			return solver, true
		}
	}
	return solver{}, false
}

func solverChoices() string {
	keys := make([]string, 0, len(solvers))
	for _, solver := range solvers {
		keys = append(keys, solver.key)
	}
	return strings.Join(keys, ", ")
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

// readBoardInput returns the board supplied to -board. A dash reads all of
// standard input; an omitted flag keeps the interactive/default-board flow.
func readBoardInput(boardFlag string) (string, error) {
	if boardFlag == "-" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read board from standard input: %w", err)
		}
		return string(input), nil
	}
	if boardFlag == "" {
		return chooseBoard(), nil
	}
	return boardFlag, nil
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
	return sudoku.NewBoard(nums)
}
