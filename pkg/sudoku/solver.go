package sudoku

import "math/bits"

type SudokuSolver func(*Board) bool

func NaiveBacktrackSolve(b *Board) bool {
	if b.IsCompleted() {
		return true
	}

	// Traverse to the first empty cell found
	for i := range 9 {
		for j := range 9 {
			if b[i][j] != Empty {
				continue
			}

			// For each possible digit
			for n := One; n <= Nine; n++ {
				b[i][j] = n
				if b.IsValid() {
					solved := NaiveBacktrackSolve(b)
					if solved {
						return true
					}
				}
				b[i][j] = Empty
			}

			// No Solutions Found
			return false
		}
	}

	panic("Invalid Sudoku Given")
}

// constraints tracks, as bitmasks, which digits are already used in each row,
// column, and 3x3 box. Bit d (1..9) is set when digit d is present. This lets
// the solver test and update placements in O(1) instead of rescanning the board.
type constraints struct {
	rows, cols, boxes [9]uint16
}

// boxIndex maps a cell coordinate to its 0..8 box number.
func boxIndex(i, j int) int {
	return (i/3)*3 + j/3
}

// newConstraints builds the bitmask state from the given board. It returns
// false if the board's existing clues already violate Sudoku rules.
func newConstraints(b *Board) (constraints, bool) {
	var c constraints
	for i := range 9 {
		for j := range 9 {
			d := b[i][j]
			if d == Empty {
				continue
			}

			bit := uint16(1) << d
			box := boxIndex(i, j)
			if c.rows[i]&bit != 0 || c.cols[j]&bit != 0 || c.boxes[box]&bit != 0 {
				return c, false
			}
			c.rows[i] |= bit
			c.cols[j] |= bit
			c.boxes[box] |= bit
		}
	}
	return c, true
}

// BacktrackSolve solves the board using backtracking with incremental
// constraint tracking. Candidate placements are validated against row/column/box
// bitmasks in O(1), avoiding the full-board rescan the naive solver performs.
func BacktrackSolve(b *Board) bool {
	c, ok := newConstraints(b)
	if !ok {
		return false
	}
	return backtrack(b, &c)
}

func backtrack(b *Board, c *constraints) bool {
	// Find the first empty cell.
	row, col := -1, -1
	for i := range 9 {
		for j := range 9 {
			if b[i][j] == Empty {
				row, col = i, j
				break
			}
		}
		if row != -1 {
			break
		}
	}

	// No empty cell remains: every placement was constraint-checked, so the
	// board is a valid complete solution.
	if row == -1 {
		return true
	}

	box := boxIndex(row, col)
	used := c.rows[row] | c.cols[col] | c.boxes[box]
	for n := One; n <= Nine; n++ {
		bit := uint16(1) << n
		if used&bit != 0 {
			continue
		}

		b[row][col] = n
		c.rows[row] |= bit
		c.cols[col] |= bit
		c.boxes[box] |= bit

		if backtrack(b, c) {
			return true
		}

		b[row][col] = Empty
		c.rows[row] &^= bit
		c.cols[col] &^= bit
		c.boxes[box] &^= bit
	}

	// No digit works here: backtrack.
	return false
}

// allDigits is a bitmask with bits 1..9 set (the valid Sudoku digits).
const allDigits uint16 = 0x3FE

// BacktrackMRVSolve solves the board with backtracking guided by the Minimum
// Remaining Values heuristic: at each step it branches on the empty cell with
// the fewest candidates. This prunes the search far more aggressively than
// always taking the first empty cell.
func BacktrackMRVSolve(b *Board) bool {
	c, ok := newConstraints(b)
	if !ok {
		return false
	}
	return backtrackMRV(b, &c)
}

func backtrackMRV(b *Board, c *constraints) bool {
	// Find the empty cell with the fewest candidates.
	row, col := -1, -1
	var bestCandidates uint16
	bestCount := 10 // more than any cell can have (max 9)

search:
	for i := range 9 {
		for j := range 9 {
			if b[i][j] != Empty {
				continue
			}

			candidates := allDigits &^ (c.rows[i] | c.cols[j] | c.boxes[boxIndex(i, j)])
			count := bits.OnesCount16(candidates)
			if count == 0 {
				// An empty cell with no candidates: this branch is dead.
				return false
			}
			if count < bestCount {
				bestCount, bestCandidates = count, candidates
				row, col = i, j
				if count == 1 {
					break search // can't do better than a forced cell
				}
			}
		}
	}

	// No empty cell remains: the board is a valid complete solution.
	if row == -1 {
		return true
	}

	box := boxIndex(row, col)
	for candidates := bestCandidates; candidates != 0; {
		bit := candidates & -candidates // lowest set bit
		candidates &^= bit

		b[row][col] = Digit(bits.TrailingZeros16(bit))
		c.rows[row] |= bit
		c.cols[col] |= bit
		c.boxes[box] |= bit

		if backtrackMRV(b, c) {
			return true
		}

		b[row][col] = Empty
		c.rows[row] &^= bit
		c.cols[col] &^= bit
		c.boxes[box] &^= bit
	}

	// No candidate works here: backtrack.
	return false
}
