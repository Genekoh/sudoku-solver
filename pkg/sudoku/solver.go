package sudoku

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
