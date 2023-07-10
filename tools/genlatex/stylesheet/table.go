package stylesheet

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

type Table struct {
	ColumnSpec   []*ColumnSpec `yaml:"columnSpec"`   // ColumnSpec's for each column
	HeaderPrefix string        `yaml:"HeaderPrefix"` // Content before a thead
	HeaderSuffix string        `yaml:"headerSuffix"` // Content after a thead
	Style        `yaml:",inline"`
}

type ColumnSpec struct {
	Hidden   bool    `yaml:"hidden"`             // If true the entire column is removed from output
	FontSize string  `yaml:"fontSize"`           // LaTeX font size, e.g. \normalsize \footnotesize
	ColType  string  `yaml:"colType"`            // LaTeX definition for column, e.g. l, r, c, p{width}
	ColWidth float64 `yaml:"colWidth,omitempty"` // Width of column
	ColUnit  string  `yaml:"colUnit,omitempty"`  // Unit for width
	Style    `yaml:",inline"`
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

	t.VerticalAlign = DefString(t.VerticalAlign, "c")

	// Ensure we have default's for each field
	for _, cs := range t.ColumnSpec {
		if cs.Hidden {
			// custom type to hide a column
			cs.ColType = "H"
		} else {
			cs.FontSize = DefString(cs.FontSize, "\\normalsize")
			cs.ColType = DefString(cs.ColType, "c")
		}
	}
}

// GetColumn gets the ColumnSpec for a specific column.
// If n < 0 then 0 is used.
// If n is beyond the number of specified columns then the last one is returned.
func (t *Table) GetColumn(n int) *ColumnSpec {
	if n < 0 {
		n = 0
	}

	l := len(t.ColumnSpec) - 1
	if n > l {
		n = l
	}

	return t.ColumnSpec[n]
}

func (t *Table) GetColumnWidth(col, span int) float64 {
	width := 0.0
	for ; span > 0; col++ {
		cd := t.GetColumn(col)
		if !cd.Hidden {
			// If column is not hidden then include it and move to next one
			width = width + cd.ColWidth
			span--
		} else if col >= (len(t.ColumnSpec) - 1) {
			// If last column is hidden then stop here
			return width
		}
	}
	return width
}

func (t *Table) GetColumnDef(col, span int) string {
	cDef := t.GetColumn(col)
	unit := cDef.ColUnit
	def := cDef.ColType

	width := FloatToString(t.GetColumnWidth(col, span), unit)
	if width != "" {
		return def + "{" + width + "}"
	}

	return def
}

// FloatToString converts float64 to string but removes unnecessary trailing 0's
func FloatToString(f float64, u string) string {
	if f <= 0.0001 {
		return ""
	}

	s := strconv.FormatFloat(f, 'f', -1, 64)
	if strings.Contains(s, ".") && !strings.HasSuffix(s, ".0") {
		s = strings.TrimRight(s, "0")
	}
	return s + u
}

func (t *Table) GetColDefs(cols int) string {
	if cols < 1 {
		cols = 1
	}

	var defs []string
	// Note we want cols entries in defs, not col here because
	// we need to skip hidden columns
	for col := 0; len(defs) < cols; col++ {
		if !t.GetColumn(col).Hidden {
			defs = append(defs, t.GetColumnDef(col, 1))
		}
	}
	return strings.TrimSpace(strings.Join(defs, " "))
}
