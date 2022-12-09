package table

import (
	"github.com/peter-mount/documentation/tools/latex/parser"
	"github.com/peter-mount/documentation/tools/latex/util"
	"golang.org/x/net/html"
	"strings"
)

type Table struct {
	w         *util.Writer
	n         *html.Node
	caption   *html.Node
	maxCols   int
	rows      int
	inTr      bool
	cellCount int
	data      [][]*Cell
	curRow    []*Cell
}

func (t *Table) Parse(n *html.Node) error {
	t.n = n
	err := parser.Traverse(n, html.ElementNode, t.parse)
	if err != nil {
		return err
	}
	t.fixStructure()
	return nil
}

func (t *Table) parse(n *html.Node) error {
	nn := n.DataAtom.String()
	switch nn {
	case "caption":
		t.caption = n
	case "tr":
		t.inTr = true
		t.cellCount = 0
		t.curRow = nil
		_ = parser.TraverseChildren(n, 0, t.parse)
		if t.cellCount > t.maxCols {
			t.maxCols = t.cellCount
		}
		t.inTr = false
		t.rows++
		t.data = append(t.data, t.curRow)
	case "th", "td":
		c := &Cell{
			rowspan: util.GetAttributeIntDefault(n, "rowspan", 1),
			colspan: util.GetAttributeIntDefault(n, "colspan", 1),
			n:       n,
			header:  nn == "th",
		}
		t.curRow = append(t.curRow, c)
		if c.header {
			t.cellCount += c.colspan
		} else {
			t.cellCount++
		}
	}
	return nil
}

func (t *Table) Write(w *util.Writer) error {
	t.w = w

	w.WriteString("\n%% Table %d rows %d cols\n", t.rows, t.maxCols).
		WriteString("\\begin{table}{\n").
		WriteString("\\tiny\n\\begin{tabular}{|%s}\n", strings.Repeat("c|", t.maxCols))

	for _, row := range t.data {
		//t.w.WriteString("\\hline\n")

		for i, cell := range row {
			// Cell separator if not the first one
			if i > 0 {
				t.w.WriteString(" & ")
			}

			err := t.writeCell(cell)
			if err != nil {
				return err
			}
		}

		// Row terminator
		t.w.WriteString(" \\\\\n")
		//t.w.WriteString(" \\\\ [0.5ex]\n")
	}

	w.WriteString("\\hline\n\\end{tabular}\n}\n")
	if t.caption != nil {
		w.WriteString("\\caption{")
		_ = parser.TraverseChildren(t.caption, 0, w.WriteHtml)
		w.WriteString("}\n")
	}
	w.WriteString("\\end{table}\n")

	return nil
}
