# sudoku-solver

A Go Sudoku solver with a small command-line interface and several solving
algorithms for comparing different approaches.

## Requirements

- Go 1.22.0 or newer

## Usage

Run the CLI with:

```sh
go run .
```

With no flags, the program prompts you to choose a solver and then asks for a
starting board. Press Enter at the board prompt to use the built-in sample
puzzle.

You can also pass everything non-interactively:

```sh
go run . -solver mrv -board 210903000000071800000002006007000013600000057001000000062050000000600400304000000
```

The board must contain 81 cells in row-major order. Digits `1` through `9` are
givens, and `0` or `.` represents an empty cell. Whitespace is ignored, so a
multi-line board is valid.

Read a board from standard input with `-board -`:

```sh
printf "210903000000071800000002006007000013600000057001000000062050000000600400304000000" | go run . -solver dlx -board -
```

Show available flags:

```sh
go run . -h
```

## Solvers

Choose a solver with `-solver`:

| Key | Name | Description |
| --- | --- | --- |
| `naive` | Naive backtracking | Tries digits in the first empty cell and rescans the board for validity after each placement. |
| `backtrack` | Backtracking | Uses row, column, and box bitmasks to check candidate placements in constant time. |
| `mrv` | MRV backtracking | Uses the same bitmask constraints, plus the Minimum Remaining Values heuristic to choose the empty cell with the fewest candidates. |
| `dlx` | Dancing links | Models Sudoku as an exact cover problem and solves it with Knuth's Algorithm X using Dancing Links. |

## Development

Run the test suite:

```sh
go test ./...
```

Run benchmarks:

```sh
go test -bench=. ./pkg/sudoku
```
