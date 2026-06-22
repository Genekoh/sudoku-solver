package sudoku

func isValidSudokuSet(set [9]Digit) bool {
	// seen is a bitmask where bit d marks digit d as already present.
	var seen uint16
	for _, c := range set {
		if c == Empty {
			continue
		}
		if c < One || c > Nine {
			return false
		}

		bit := uint16(1) << c
		if seen&bit != 0 {
			return false
		}
		seen |= bit
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
		col := [9]Digit{}
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
			block := [9]Digit{}
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
			if c == Empty {
				return false
			}
		}
	}

	return true
}

func (b Board) IsCompleted() bool {
	return b.IsValid() && b.IsFilled()
}
