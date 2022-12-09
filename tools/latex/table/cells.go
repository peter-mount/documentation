package table

import (
	"fmt"
	"github.com/peter-mount/documentation/tools/latex/parser"
	"golang.org/x/net/html"
)

type Cell struct {
	rowspan int
	colspan int
	header  bool
	n       *html.Node
}

func (t *Table) getRow(r int) []*Cell {
	for len(t.data) <= r {
		t.data = append(t.data, nil)
	}
	return t.data[r]
}

// fixStructure ensures the table is laid out as LaTeX expects by
// filling in rowspan's correctly
func (t *Table) fixStructure() {
	for r, _ := range t.data {

		// Shift cells so row spanned cells have an empty space present
		for c, cell := range t.data[r] {
			if cell.rowspan > 1 {
				for i := 1; i < cell.rowspan; i++ {
					if (r + i) < len(t.data) {
						if c == 0 {
							t.data[r+i] = append([]*Cell{&Cell{}}, t.data[r+i]...)
						} else if c < len(t.data[r+i]) {
							t.data[r+i] = append(t.data[r+i][:c+1], t.data[r+i][c:]...)
							t.data[r+i][c] = &Cell{}
						} else {
							for c > len(t.data[r+i]) {
								t.data[r+i] = append(t.data[r+i], &Cell{})
							}
						}
					}
				}
			}
		}

		// Ensure row is full length
		for len(t.data[r]) < t.maxCols {
			t.data[r] = append(t.data[r], &Cell{})
		}
	}

	for r, row := range t.data {
		fmt.Printf("%02d|", r)
		for _, cell := range row {
			if cell.n == nil {
				fmt.Print(" |")
			} else {
				fmt.Print("*|")
			}
		}
		fmt.Println()
	}
}

func (t *Table) writeCell(c *Cell) error {
	// No node then it's an empty cell
	if c.n == nil {
		return nil
	}

	if c.colspan > 1 {
		t.w.WriteString("\\multicolumn{%d}{|c|}{", c.colspan)
	}

	if c.rowspan > 1 {
		t.w.WriteString("\\multirow{%d}{}{", c.rowspan)
	}

	if c.header {
		t.w.WriteString("\\textbf{")
	}

	err := parser.TraverseChildren(c.n, 0, t.w.WriteHtml)
	if err != nil {
		return err
	}

	if c.header {
		t.w.WriteString("}")
	}

	if c.rowspan > 1 {
		t.w.WriteString("}")
	}

	if c.colspan > 1 {
		t.w.WriteString("}")
	}

	return nil
}
