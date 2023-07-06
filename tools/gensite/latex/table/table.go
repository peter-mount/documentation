package table

import (
	"github.com/peter-mount/documentation/tools/gensite/latex/parser"
	"github.com/peter-mount/documentation/tools/gensite/latex/util"
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
	inHead    bool
	cellCount int
	head      CellSet
	body      CellSet
	curRow    []*Cell
}

func (t *Table) Parse(n *html.Node) error {
	t.n = n
	err := parser.Traverse(n, html.ElementNode, t.parse)
	if err != nil {
		return err
	}
	t.head.fixStructure(t.maxCols)
	t.body.fixStructure(t.maxCols)
	return nil
}

func (t *Table) parse(n *html.Node) error {
	nn := n.DataAtom.String()
	switch nn {
	case "caption":
		t.caption = n
	case "thead":
		t.inHead = true
	case "tbody":
		t.inHead = false
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
		if t.inHead {
			t.head = append(t.head, t.curRow)
		} else {
			t.body = append(t.body, t.curRow)
		}
	case "th", "td":
		c := &Cell{
			rowspan: util.GetAttributeIntDefault(n, "rowspan", 1),
			colspan: util.GetAttributeIntDefault(n, "colspan", 1),
			n:       n,
			header:  nn == "th",
			align:   AlignCenter,
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
		WriteString("\\tiny\n\\begin{tabularx}{\\textwidth}{|%s}\n", strings.Repeat("c|", t.maxCols))

	w.WriteString("\\hline\n")

	err := t.head.Write(w)
	if err != nil {
		return err
	}

	w.WriteString("\\hline\n")

	err = t.body.Write(w)
	if err != nil {
		return err
	}

	w.WriteString("\\hline\n\\end{tabularx}\n}\n")
	if t.caption != nil {
		w.WriteString("\\caption{")
		_ = parser.TraverseChildren(t.caption, 0, w.WriteHtml)
		w.WriteString("}\n")
	}
	w.WriteString("\\end{table}\n")

	return nil
}
