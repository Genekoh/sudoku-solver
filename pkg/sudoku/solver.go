package sudoku

type SudokuSolver func(*Board) bool

func NaiveBacktrackSolve(b *Board) bool {
	if b.IsCompleted() {
		return true
	}

	// Traverse to the first empty cell found
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
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
