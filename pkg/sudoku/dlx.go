package sudoku

// DancingLinksSolve solves the board by modeling it as an exact cover problem
// and applying Knuth's Algorithm X with the Dancing Links technique.
//
// The exact cover matrix has 324 constraint columns, in four groups of 81:
//   - cell:   each (row, col) holds exactly one digit
//   - row:    each (row, digit) appears exactly once
//   - col:    each (col, digit) appears exactly once
//   - box:    each (box, digit) appears exactly once
//
// Each matrix row represents placing a specific digit in a specific cell and
// covers exactly those four constraints. Givens are encoded by only emitting
// the single matrix row consistent with the clue, so a complete cover is a
// valid completed board.
func DancingLinksSolve(b *Board) bool {
	m := newDLX()

	for r := range 9 {
		for c := range 9 {
			given := b[r][c]
			for d := One; d <= Nine; d++ {
				if given != Empty && d != given {
					continue
				}
				m.addRow(r, c, int(d))
			}
		}
	}

	solution := make([]*dlxNode, 0, 81)
	if !m.search(&solution) {
		return false
	}

	for _, n := range solution {
		r, c, d := decodeRow(n.rowID)
		b[r][c] = Digit(d)
	}
	return true
}

// dlxNode is a node in the four-way circular linked structure. Column header
// nodes additionally use size; data nodes use rowID and point to their header
// via col.
type dlxNode struct {
	left, right, up, down *dlxNode
	col                   *dlxNode
	size                  int // column headers only: number of data nodes
	rowID                 int // data nodes only: encodes (row, col, digit)
}

type dlx struct {
	root    *dlxNode
	columns [324]*dlxNode
}

// Constraint-column group offsets.
const (
	cellOffset = 0
	rowOffset  = 81
	colOffset  = 162
	boxOffset  = 243
)

func boxOf(r, c int) int { return (r/3)*3 + c/3 }

// encodeRow maps a (row, col, digit) placement to a unique row id; digit is 1..9.
func encodeRow(r, c, d int) int { return (r*9+c)*9 + (d - 1) }

func decodeRow(id int) (r, c, d int) {
	d = id%9 + 1
	id /= 9
	c = id % 9
	r = id / 9
	return r, c, d
}

func newDLX() *dlx {
	m := &dlx{root: &dlxNode{}}
	m.root.left, m.root.right = m.root, m.root

	// Create and link 324 column headers into the header ring.
	for i := range m.columns {
		col := &dlxNode{}
		col.col = col
		col.up, col.down = col, col

		col.right = m.root
		col.left = m.root.left
		m.root.left.right = col
		m.root.left = col

		m.columns[i] = col
	}
	return m
}

// addRow appends the matrix row for placing digit d (1..9) at cell (r, c),
// wiring its four nodes into the relevant column lists.
func (m *dlx) addRow(r, c, d int) {
	cols := [4]int{
		cellOffset + r*9 + c,
		rowOffset + r*9 + (d - 1),
		colOffset + c*9 + (d - 1),
		boxOffset + boxOf(r, c)*9 + (d - 1),
	}

	id := encodeRow(r, c, d)
	var first *dlxNode
	for _, ci := range cols {
		col := m.columns[ci]
		n := &dlxNode{col: col, rowID: id}

		// Append to the bottom of the column's vertical list.
		n.down = col
		n.up = col.up
		col.up.down = n
		col.up = n
		col.size++

		// Link horizontally with the other nodes of this row.
		if first == nil {
			first = n
			n.left, n.right = n, n
		} else {
			n.right = first
			n.left = first.left
			first.left.right = n
			first.left = n
		}
	}
}

func cover(col *dlxNode) {
	col.right.left = col.left
	col.left.right = col.right
	for i := col.down; i != col; i = i.down {
		for j := i.right; j != i; j = j.right {
			j.down.up = j.up
			j.up.down = j.down
			j.col.size--
		}
	}
}

func uncover(col *dlxNode) {
	for i := col.up; i != col; i = i.up {
		for j := i.left; j != i; j = j.left {
			j.col.size++
			j.down.up = j
			j.up.down = j
		}
	}
	col.right.left = col
	col.left.right = col
}

// chooseColumn returns the uncovered column with the fewest nodes (Knuth's S
// heuristic), which minimizes branching.
func (m *dlx) chooseColumn() *dlxNode {
	best := m.root.right
	for col := best; col != m.root; col = col.right {
		if col.size < best.size {
			best = col
		}
	}
	return best
}

func (m *dlx) search(solution *[]*dlxNode) bool {
	if m.root.right == m.root {
		return true // every constraint satisfied
	}

	col := m.chooseColumn()
	if col.size == 0 {
		return false // a constraint with no remaining option: dead end
	}

	cover(col)
	for r := col.down; r != col; r = r.down {
		*solution = append(*solution, r)
		for j := r.right; j != r; j = j.right {
			cover(j.col)
		}

		if m.search(solution) {
			return true
		}

		// Backtrack.
		*solution = (*solution)[:len(*solution)-1]
		for j := r.left; j != r; j = j.left {
			uncover(j.col)
		}
	}
	uncover(col)
	return false
}
