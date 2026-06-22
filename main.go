package main

import (
	"fmt"

	"github.com/Genekoh/sudoku-solver/pkg/sudoku"
)

func main() {
	b := sudoku.NewBoard([81]int{
		2, 1, 0, 9, 0, 3, 0, 0, 0,
		0, 0, 0, 0, 7, 1, 8, 0, 0,
		0, 0, 0, 0, 0, 2, 0, 0, 6,
		0, 0, 7, 0, 0, 0, 0, 1, 3,
		6, 0, 0, 0, 0, 0, 0, 5, 7,
		0, 0, 1, 0, 0, 0, 0, 0, 0,
		0, 6, 2, 0, 5, 0, 0, 0, 0,
		0, 0, 0, 6, 0, 0, 4, 0, 0,
		3, 0, 4, 0, 0, 0, 0, 0, 0,
	})

	fmt.Println("Given:")
	sudoku.PrintBoard(*b)
	fmt.Println("\nSolved:")

	if !sudoku.NaiveBacktrackSolve(b) {
		fmt.Printf("Failed to solve")
	} else {
		sudoku.PrintBoard(*b)
	}
}
