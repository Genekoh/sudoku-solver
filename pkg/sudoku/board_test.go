package sudoku_test

import (
	"strings"
	"testing"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

func TestNewBoardRejectsDigitsOutsideRange(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{name: "negative", value: -1},
		{name: "too-large", value: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cells [81]int
			cells[17] = tt.value

			if _, err := sudoku.NewBoard(cells); err == nil {
				t.Fatalf("NewBoard accepted out-of-range digit %d", tt.value)
			}
		})
	}
}

func TestIsValidRejectsOutOfRangeDigit(t *testing.T) {
	var board sudoku.Board
	board[0][0] = sudoku.Digit(10)

	if board.IsValid() {
		t.Fatal("IsValid accepted a digit outside 1..9")
	}
}

func TestBoardString(t *testing.T) {
	board, err := sudoku.NewBoard([81]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if err != nil {
		t.Fatalf("NewBoard: %v", err)
	}

	lines := strings.Split(board.String(), "\n")
	if len(lines) != 11 {
		t.Fatalf("String produced %d lines, want 11", len(lines))
	}
	if got, want := lines[0], "1 2 3 | 4 5 6 | 7 8 9"; got != want {
		t.Errorf("first row = %q, want %q", got, want)
	}
	if got, want := lines[1], "x x x | x x x | x x x"; got != want {
		t.Errorf("empty row = %q, want %q", got, want)
	}
}
