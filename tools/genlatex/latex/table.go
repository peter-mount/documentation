package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"golang.org/x/net/html"
	"math"
	"strings"
)

// Return the max number of columns in a table n
func getTableColCount(n *html.Node, ctx context.Context) int {
	cols := 0
	rowCols := 0
	var f parser.Handler
	f = func(n *html.Node, ctx context.Context) error {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "thead", "tbody":
				return f.HandleChildren(n, ctx)

			case "tr":
				rowCols = 0
				_ = f.HandleChildren(n, ctx)
				if rowCols > cols {
					cols = rowCols
				}
				return nil

			case "th", "td":
				colSpan, _ := parser.GetAttrInt(n, "colspan", 1)
				rowCols = rowCols + colSpan
				return nil
			}
		}
		return f.HandleChildren(n, ctx)
	}
	_ = f.HandleChildren(n, ctx)
	return cols
}

// tableState used to keep track of which columns are in use
type tableState struct {
	// Number of columns
	cols int
	// Column inactive state - set if a column is in a rowspan
	// 0 means it's active and expects a value.
	// Hidden columns get this set to MaxInt
	colRowSpan []int
	// Column definition cache
	colDefs []*stylesheet.ColumnSpec
}

func newTableState(cols int, table *stylesheet.Table, ctx context.Context) context.Context {
	ts := &tableState{
		cols:    cols,
		colDefs: table.ColumnSpec,
	}
	ts.reset()
	return context.WithValue(ctx, tableStateKey, ts)
}

func tableStateFromContext(ctx context.Context) *tableState {
	v := ctx.Value(tableStateKey)
	if v != nil {
		return v.(*tableState)
	}
	return nil
}

// reset() called at end of thead and tbody
func (ts *tableState) reset() {
	// init colRowSpan so hidden columns are ignored, all others are initially 0
	for col := 0; col < ts.cols; col++ {
		rowSpan := 0
		if ts.getCol(col).Hidden {
			rowSpan = math.MaxInt
		}
		ts.colRowSpan = append(ts.colRowSpan, rowSpan)
	}
}

func (ts *tableState) getCol(col int) *stylesheet.ColumnSpec {
	l := len(ts.colDefs)
	if col >= l {
		col = l - 1
	}
	return ts.colDefs[col]
}

// adjustCellNumber increments cell if it's on a column that's part of a row span
func (ts *tableState) adjustCellNumber(col int) int {
	for ; col < ts.cols; col++ {
		if ts.colRowSpan[col] > 0 {
			ts.colRowSpan[col]--
		}
	}
	return col
}

// incrementCellNumber adds colSpan to cell accounting for any hidden Columns
func (ts *tableState) incrementCellNumber(col, colSpan int) int {
	for ; colSpan > 0; colSpan-- {
		// If hidden then skip this column
		//for col < ts.cols && ts.getCol(col).Hidden {
		//	col++
		//}

		// Increment col
		col++
	}
	return col
}

// setRowSpan marks a column has a row span
func (ts *tableState) setRowSpan(col, rowSpan int) {
	// -1 because we are already in this column
	ts.colRowSpan[col] = rowSpan - 1
}

func (c *Converter) table(n *html.Node, ctx context.Context) error {
	// If we are nested then add a \\ prefix
	if stylesheet.TableFromContext(ctx) != nil {
		_ = Write(ctx, ' ', '\\', '\\', '\n')
	}

	// Needed to save state otherwise headers spanning pages will break subsequent tables
	return group(c.tableImpl, n, ctx)
}

func (c *Converter) tableImpl(n *html.Node, ctx context.Context) error {

	// if nested then tabular otherwise the outer table is a longtable
	tType := "longtable"
	if stylesheet.TableFromContext(ctx) != nil {
		tType = "tabular"

		// Reset header flag incase we are also in the thead section of the parent
		ctx = context.WithValue(ctx, insideTableHeaderKey, nil)
	}

	// how many columns are we?
	cols := getTableColCount(n, ctx)

	style := c.Stylesheet().GetStyle(n)
	table := c.Stylesheet().GetTable(n)
	ctx = table.WithContext(ctx)

	// Initialise table state
	ctx = newTableState(cols, table, ctx)

	// Note: @{} either side of the col specifiers tells LaTeX not to add inter-column spacing
	// before and after the first & last columns respectively. Without that it would be
	// wasted space on the page.
	//
	// [t] means vertical align cells to the top. Defaults to c otherwise, other option is b
	err := Writef(ctx,
		"\\begin{%s}[%s]{@{} %s @{}}\n",
		tType,
		stylesheet.DefStrings(
			style.VerticalAlign,
			table.VerticalAlign,
		),
		table.GetColDefs(cols),
	)

	if err == nil {
		var f parser.Handler
		f = func(n *html.Node, ctx context.Context) error {
			if n.Type == html.ElementNode {
				switch n.Data {
				case "thead":
					return c.tableHeader(f, n, ctx)

				case "tbody":
					return parser.HandleChildren(f, n, ctx)

				case "tr":
					return c.tr(n, ctx)

				}
			}
			return f.HandleChildren(n, ctx)
		}
		err = f.HandleChildren(n, ctx)
	}

	if err == nil {
		err = Writef(ctx, "\\end{%s}\n", tType)
	}

	return err
}

func (c *Converter) tableHeader(f parser.Handler, n *html.Node, ctx context.Context) error {
	// Mark we are in a header
	ctx = context.WithValue(ctx, insideTableHeaderKey, true)

	table := stylesheet.TableFromContext(ctx)
	var err error

	if table.HeaderPrefix != "" {
		err = WriteStringLn(ctx, table.HeaderPrefix)
	}

	if err == nil {
		err = parser.HandleChildren(f, n, ctx)
	}

	if err == nil && table.HeaderSuffix != "" {
		err = WriteStringLn(ctx, table.HeaderSuffix)
	}

	if err == nil {
		err = WriteStringLn(ctx, `\endhead`)
	}

	return err
}

func (c *Converter) tableCaption(n *html.Node, ctx context.Context) error {
	return nil // TODO handleChildren(n, ctx)
}

const (
	// Marker to indicate we are inside the header
	insideTableHeaderKey = "inside.table.header"
	tableStateKey        = "table.state"
)

func insideTableHeader(ctx context.Context) bool {
	return ctx.Value(insideTableHeaderKey) != nil
}

func (c *Converter) tr(n *html.Node, ctx context.Context) error {

	// The table definition
	table := stylesheet.TableFromContext(ctx)
	ts := tableStateFromContext(ctx)

	cell := 0
	firstCell := true
	err := parser.HandleChildren(func(n *html.Node, ctx context.Context) (err error) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "th", "td":

				style := c.Stylesheet().GetStyle(n)
				colSpan, _ := parser.GetAttrInt(n, "colspan", 1)
				rowSpan, _ := parser.GetAttrInt(n, "rowspan", 1)
				multiCol, multiRow := colSpan > 1, rowSpan > 1

				// Adjust cell to skip columns within a row span

				// Handle cell separator
				if firstCell {
					firstCell = false
				} else {
					err = WriteString(ctx, " & ")
				}

				// Look ahead, if a table exists then we need to wrap the cell to allow line break before it.
				// nonTabularContent is used to set tabularCell so if table is first then it's not tabular
				// unless there's multiple tables.
				//
				// We also set tabularCell if we have a br element, so it will linebreak correctly
				tabularCell, nonTabularContent := false, false
				_ = parser.HandleChildren(func(n *html.Node, ctx context.Context) error {
					if parser.IsElement(n, "table") || parser.IsElement(n, "br") {
						tabularCell = nonTabularContent
					} else if n.Type == html.TextNode {
						// Ignore whitespace only text nodes
						nonTabularContent = strings.Trim(n.Data, " \n\t") != ""
						//} else {
						//	nonTabularContent = true
					}
					return nil
				}, n, ctx)

				if err == nil && multiCol {
					err = Writef(ctx,
						`\multicolumn{%d}{%s}{`,
						colSpan,
						// Without this we get Package array Error: Empty preamble: `l' used.
						stylesheet.DefString(table.GetColumnDef(cell, colSpan), "l"),
					)
				}
				if err == nil && multiRow {
					err = Writef(ctx, `\multirow{%d}{%s}{`,
						rowSpan,
						// width or "*" for contents natural width
						stylesheet.DefString(
							stylesheet.FloatToString(
								table.GetColumnWidth(cell, colSpan),
								table.GetColumn(cell).ColUnit),
							"*"),
					)
				}

				if err == nil && tabularCell {
					// Wrap cell content within a tabular block, so we can then use \\ as a line break
					err = Writef(ctx,
						`\begin{tabular}[%s]{@{}l@{}}`,
						stylesheet.DefStrings(
							style.VerticalAlign,
							table.GetColumn(cell).VerticalAlign,
							table.VerticalAlign),
					)
				}

				if err == nil {
					err = handleChildren(n, ctx)
				}

				if err == nil && tabularCell {
					err = WriteString(ctx, `\end{tabular}`)
				}
				if err == nil && multiRow {
					err = Write(ctx, '}')
				}
				if err == nil && multiCol {
					err = Write(ctx, '}')
				}

				// Increment cell number
				cell = ts.incrementCellNumber(cell, colSpan)
			default:
			}
		}
		return
	}, n, ctx)

	// End table row marker
	if err == nil {
		err = WriteStringLn(ctx, ` \\`)
	}

	return err
}
