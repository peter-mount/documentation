package latex

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"strings"
)

type Table struct {
	Type    string            // Type, e.g. "longtblr" or "tblr"
	Table   *stylesheet.Table // Stylesheet
	Rows    []*TableRow       // Table rows
	Caption string            // Table caption, if any
}

func (t *Table) Finalise() {
	// Delete hidden columns - must be in reverse order
	// otherwise you can end up deleting the wrong columns
	for col := len(t.Table.ColumnSpec) - 1; col >= 0; col-- {
		cs := t.Table.ColumnSpec[col]
		if cs.Hidden {
			t.DeleteColumn(col)
		}
	}
}

func (t *Table) DeleteColumn(col int) {
	for _, r := range t.Rows {
		r.DeleteColumn(col)
	}
}

func (t *Table) String() string {
	// Work out max number of columns
	numCols := 0
	for _, r := range t.Rows {
		nc := len(r.Cells)
		if nc > numCols {
			numCols = nc
		}
	}

	// Ensure all rows have the same number of columns
	for _, r := range t.Rows {
		r.fixCells(numCols - 1)
	}

	var p []string
	if t.Caption != "" {
		p = append(p, "caption={", t.Caption, "},")
	}

	// Generate the table
	var s []string

	s = append(s, "\\begin{", t.Type, "}")

	if len(p) > 0 {
		s = append(s, "[")
		s = append(s, p...)
		s = append(s, "]")
	}

	s = append(s, "{@{} ", t.Table.GetColDefs(numCols), " @{}}")

	for _, r := range t.Rows {
		s = r.Append(s)
	}

	s = append(s, fmt.Sprintf("\\end{%s}\n", t.Type))

	return strings.Join(s, "")
}

func (t *Table) RowCount() int {
	return len(t.Rows)
}

func (t *Table) SetCell(r, c int, cell *TableCell) {
	// Ensure we have enough rows
	for r >= len(t.Rows) {
		t.Rows = append(t.Rows, &TableRow{})
	}
	t.Rows[r].SetCell(c, cell)
	//fmt.Printf("setCell %d, %d\n", r, c)
}

type TableRow struct {
	Cells []*TableCell
}

func (t *TableRow) DeleteColumn(col int) {
	if col < len(t.Cells) {
		a := make([]*TableCell, 0)
		a = append(a, t.Cells[:col]...)
		a = append(a, t.Cells[col+1:]...)
		t.Cells = a
	}
}

func (t *TableRow) Append(s []string) []string {
	for i, c := range t.Cells {
		if i > 0 {
			s = append(s, " & ")
		}
		s = c.Append(s)
	}
	return append(s, " \\\\\n")
}

func (t *TableRow) fixCells(c int) {
	// Ensure we have enough cells in the row
	for c >= len(t.Cells) {
		t.Cells = append(t.Cells, nil)
	}
}

func (t *TableRow) SetCell(c int, cell *TableCell) {
	if cell.ColSpan < 1 {
		cell.ColSpan = 1
	}
	if cell.RowSpan < 1 {
		cell.RowSpan = 1
	}
	t.fixCells(c + cell.ColSpan - 1)
	t.Cells[c] = cell
}

type TableCell struct {
	ColSpan int
	RowSpan int
	Text    string
}

func (c *TableCell) Append(s []string) []string {
	if c != nil {
		if c.RowSpan > 1 || c.ColSpan > 1 {
			var a []string
			if c.ColSpan > 1 {
				a = append(a, fmt.Sprintf("c=%d", c.ColSpan))
			}
			if c.RowSpan > 1 {
				a = append(a, fmt.Sprintf("r=%d", c.RowSpan))
			}
			s = append(s, `\SetCell[`, strings.Join(a, ","), `]{c}`)
		}
		s = append(s, c.Text)
	}

	return s
}
