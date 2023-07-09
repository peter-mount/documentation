package latex

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"golang.org/x/net/html"
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

func cellDelimiter(firstCell bool, ctx context.Context) (bool, error) {
	if firstCell {
		return false, nil
	} else {
		return false, WriteString(ctx, " & ")
	}
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

				// Handle cell separators
				cell, firstCell, err = ts.addEmptyCells(cell, firstCell, ctx)
				firstCell, err = cellDelimiter(firstCell, ctx)

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

				// content to add before and after a cell's content has been completed
				cellPrefix, cellSuffix := "", ""

				// Handle the 4 kinds of cell, single (do nothing), multiCol, multiRow and both multi Row & Col
				if err == nil {
					switch {
					case multiCol && multiRow:
						ts.setRowSpan(cell, rowSpan)

						cellPrefix = fmt.Sprintf(
							`\multicolumn{%d}{%s}{\multirow{%d}{%s}{`,
							colSpan,
							// For now these are all top-left
							"tl",
							rowSpan,
							// width or "*" for contents natural width
							stylesheet.DefString(
								stylesheet.FloatToString(
									table.GetColumnWidth(cell, colSpan),
									table.GetColumn(cell).ColUnit),
								"*"),
						)

						cellSuffix = "}}"

					case multiCol:
						cellPrefix = fmt.Sprintf(
							`\multicolumn{%d}{%s}{`,
							colSpan,
							// Without this we get Package array Error: Empty preamble: `l' used.
							stylesheet.DefString(table.GetColumnDef(cell, colSpan), "l"),
						)

						cellSuffix = "}"

					case multiRow:
						ts.setRowSpan(cell, rowSpan)

						cellPrefix = fmt.Sprintf(
							`\multirow{%d}{%s}{`,
							rowSpan,
							// width or "*" for contents natural width
							stylesheet.DefString(
								stylesheet.FloatToString(
									table.GetColumnWidth(cell, colSpan),
									table.GetColumn(cell).ColUnit),
								"*"),
						)

						cellSuffix = "}"
					}
				}

				if err == nil && tabularCell {
					// Wrap cell content within a tabular block, so we can then use \\ as a line break
					cellPrefix = cellPrefix + fmt.Sprintf(
						`\begin{tabular}[%s]{@{}l@{}}`,
						stylesheet.DefStrings(
							style.VerticalAlign,
							table.GetColumn(cell).VerticalAlign,
							table.VerticalAlign),
					)
					cellSuffix = `\end{tabular}` + cellSuffix
				}

				if err == nil && cellPrefix != "" {
					err = WriteString(ctx, cellPrefix)
				}

				if err == nil {
					err = handleChildren(n, ctx)
				}

				if err == nil && cellSuffix != "" {
					err = WriteString(ctx, cellSuffix)
				}

				// Increment cell number
				cell = ts.incrementCellNumber(cell, colSpan)
			default:
			}
		}
		return
	}, n, ctx)

	if err == nil {
		// Ensure we have enough cell entries in this row
		// for when the html table is short of cells
		_, _, err = ts.addEmptyCells(cell, firstCell, ctx)
	}

	// End table row marker
	if err == nil {
		err = WriteStringLn(ctx, ` \\`)
	}

	return err
}
