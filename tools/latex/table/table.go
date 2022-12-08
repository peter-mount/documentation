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
		t.cellCount++
	}
	return nil
}

func (t *Table) Write(w *util.Writer) error {
	t.w = w

	nw := t.w.
		Begin("center").
		Begin("tabular").
		WriteString("{|| %s||}", strings.Repeat("c ", t.maxCols))

	err := parser.TraverseChildren(t.n, 0, t.write)
	if err != nil {
		return err
	}

	nw.WriteString("\\hline\n").
		End().
		End()

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
		switch n.DataAtom.String() {
		case "tr":
			t.w.WriteString("\\hline\n")
			t.cellCount = 0
			err := parser.TraverseChildren(n, 0, t.write)
			if err != nil {
				return err
			}
			t.w.WriteString(" \\\\ [0.5ex]")
			return parser.StopChildTraverse()

		case "td", "th":
			err := parser.TraverseChildren(n, 0, t.write)
			if err != nil {
				return err
			}
			t.cellCount++
			if t.cellCount < t.maxCols {
				t.w.WriteString(" & ")
			}
			return parser.StopChildTraverse()

		}
	}

	return nil
}
