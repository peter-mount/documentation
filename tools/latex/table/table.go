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
	maxCols   int
	rows      int
	inTr      bool
	cellCount int
}

func (t *Table) Parse(n *html.Node) error {
	t.n = n
	return parser.Traverse(n, html.ElementNode, t.parse)
}

func (t *Table) parse(n *html.Node) error {
	switch n.DataAtom.String() {
	case "tr":
		t.inTr = true
		t.cellCount = 0
		_ = parser.TraverseChildren(n, 0, t.parse)
		if t.cellCount > t.maxCols {
			t.maxCols = t.cellCount
		}
		t.inTr = false
		t.rows++
	case "th", "td":
		t.cellCount += util.GetAttributeIntDefault(n, "colspan", 1)
	}
	return nil
}

func (t *Table) Write(w *util.Writer) error {
	t.w = w

	w.WriteString("%% Table %d rows %d cols\n", t.rows, t.maxCols).
		WriteString("\\begin{center}\n\\begin{tabular}{|%s}\n", strings.Repeat("c|", t.maxCols))

	err := parser.TraverseChildren(t.n, 0, t.write)
	if err != nil {
		return err
	}

	w.WriteString("\\hline\n\\end{tabular}\n\\end{center}\n")

	return nil
}

func (t *Table) write(n *html.Node) error {
	switch n.Type {
	case html.TextNode:
		s := strings.TrimSpace(n.Data)
		if len(s) > 0 {
			_ = t.w.WriteString(util.EscapeText(s))
		}

	case html.ElementNode:
		en := n.DataAtom.String()
		switch en {
		case "tr":
			t.w.WriteString("\\hline\n")
			t.cellCount = 0
			err := parser.TraverseChildren(n, 0, t.write)
			if err != nil {
				return err
			}
			t.w.WriteString(" \\\\ [0.5ex]\n")
			return parser.StopChildTraverse()

		case "td", "th":
			if t.cellCount > 0 {
				t.w.WriteString(" & ")
			}

			rows := util.GetAttributeIntDefault(n, "rowspan", 1)
			if rows > 1 {
				t.w.WriteString("\\multirow{%d}{}{", rows)
			}

			cols := util.GetAttributeIntDefault(n, "colspan", 1)
			if cols > 1 {
				t.w.WriteString("\\multicolumn{%d}{|c|}{", cols)
			}
			t.cellCount += cols

			header := en == "th"
			if header {
				t.w.WriteString("\\textbf{")
			}

			err := parser.TraverseChildren(n, 0, t.write)
			if err != nil {
				return err
			}

			if header {
				t.w.WriteString("}")
			}

			if cols > 1 {
				t.w.WriteString("}")
			}

			if rows > 1 {
				t.w.WriteString("}")
			}

			return parser.StopChildTraverse()

		}
	}

	return nil
}
