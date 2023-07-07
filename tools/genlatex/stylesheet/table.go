package stylesheet

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

type Table struct {
	ColumnSpec    []*ColumnSpec `yaml:"columnSpec"`
	VerticalAlign string        `yaml:"verticalAlign"` // Vertical Align, default c, can be t or b
}

type ColumnSpec struct {
	FontSize   string `yaml:"fontSize"`   // LaTeX font size, e.g. \normalsize \footnotesize
	Definition string `yaml:"definition"` // LaTeX definition for column, e.g. l, r, c, p{width}
}

// GetTable gets the Table for this specific node.
func (s *Stylesheet) GetTable(n *html.Node) *Table {
	for _, c := range parser.GetClass(n) {
		if t := s.Table[c]; t != nil {
			return t
		}
	}
	return s.defaultTable
}

func TableFromContext(ctx context.Context) *Table {
	v := ctx.Value(tableKey)
	if v == nil {
		return nil
	}
	return v.(*Table)
}

func (t *Table) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, tableKey, t)
}

func (t *Table) init() {
	// If none defined then create a default
	if len(t.ColumnSpec) == 0 {
		t.ColumnSpec = append(t.ColumnSpec, &ColumnSpec{})
	}

	t.VerticalAlign = defString(t.VerticalAlign, "c")

	// Ensure we have default's for each field
	for _, cs := range t.ColumnSpec {
		cs.FontSize = defString(cs.FontSize, "\\normalsize")
		cs.Definition = defString(cs.Definition, "c")
	}
}

// GetColumn gets the ColumnSpec for a specific column.
// If n < 0 then 0 is used.
// If n is beyond the number of specified columns then the last one is returned.
func (t *Table) GetColumn(n int) *ColumnSpec {
	l := len(t.ColumnSpec) - 1
	if n < 0 {
		n = 0
	}
	if n > l {
		n = l
	}
	return t.ColumnSpec[n]
}

func (t *Table) GetColDefs(cols int) string {
	if cols < 1 {
		cols = 1
	}
	var defs []string
	for col := 0; col < cols; col++ {
		defs = append(defs, t.GetColumn(col).Definition)
	}
	return strings.TrimSpace(strings.Join(defs, " "))
}
