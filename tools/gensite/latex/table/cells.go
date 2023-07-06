package table

import (
	"github.com/peter-mount/documentation/tools/gensite/latex/parser"
	"github.com/peter-mount/documentation/tools/gensite/latex/util"
	"golang.org/x/net/html"
)

type CellSet [][]*Cell

type Cell struct {
	rowspan int
	colspan int
	header  bool
	align   Align
	n       *html.Node
}

/*
	func (t *Table) getRow(r int) []*Cell {
		for len(t.body) <= r {
			t.body = append(t.body, nil)
		}
		return t.body[r]
	}
*/
func colCount(r []*Cell) int {
	c := 0
	for _, cell := range r {
		c = c + cell.colspan
	}
	return c
}

// fixStructure ensures the table is laid out as LaTeX expects by
// filling in rowspan's correctly
func (t *CellSet) fixStructure(maxCols int) {
	for r, _ := range *t {

		// Shift cells so row spanned cells have an empty space present
		for c, cell := range (*t)[r] {
			if cell.rowspan > 1 {
				for i := 1; i < cell.rowspan; i++ {
					if (r + i) < len(*t) {
						// copy colspan as multirow & multicol needs col passed down & has to be align left to work
						nc := &Cell{colspan: cell.colspan, align: AlignLeft}
						if c == 0 {
							(*t)[r+i] = append([]*Cell{nc}, (*t)[r+i]...)
						} else if c < len((*t)[r+i]) {
							(*t)[r+i] = append((*t)[r+i][:c+1], (*t)[r+i][c:]...)
							(*t)[r+i][c] = nc
						} else {
							for c > len((*t)[r+i]) {
								(*t)[r+i] = append((*t)[r+i], nc)
							}
						}
					}
				}
			}
		}

		// Ensure row is full length, accounting for col-spans by appending empty 1x1 cells
		for colCount((*t)[r]) < maxCols {
			(*t)[r] = append((*t)[r], &Cell{colspan: 1, rowspan: 1, align: AlignLeft})
		}
	}
}

func (t *CellSet) Write(w *util.Writer) error {
	for _, row := range *t {
		//t.w.WriteString("\\hline\n")

		for i, cell := range row {
			// Cell separator if not the first one
			if i > 0 {
				w.WriteString(" & ")
			}

			err := cell.Write(w)
			if err != nil {
				return err
			}
		}

		// Row terminator
		w.WriteString(" \\\\\n")
		//t.w.WriteString(" \\\\ [0.5ex]\n")
	}

	return nil
}

func (c *Cell) Write(w *util.Writer) error {
	if c.colspan > 1 {
		w.WriteString("\\multicolumn{%d}{%s}{", c.colspan, c.align)
	}

	if c.rowspan > 1 {
		w.WriteString("\\multirow{%d}{*}{", c.rowspan)
	}

	if c.n != nil {
		if c.header {
			w.WriteString("\\textbf{")
		}

		err := parser.TraverseChildren(c.n, 0, w.WriteHtml)
		if err != nil {
			return err
		}

		if c.header {
			w.WriteString("}")
		}
	}

	if c.rowspan > 1 {
		w.WriteString("}")
	}

	if c.colspan > 1 {
		w.WriteString("}")
	}

	return nil
}
