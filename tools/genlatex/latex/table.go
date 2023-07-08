package latex

import (
	"context"
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
	}

	// how many columns are we?
	cols := getTableColCount(n, ctx)

	style := c.Stylesheet().GetStyle(n)
	table := c.Stylesheet().GetTable(n)
	ctx = table.WithContext(ctx)

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
				// thead is the same as tbody just we set a flag for handling multi-page tables
				case "thead":
					return parser.HandleChildren(f, n, context.WithValue(ctx, insideTableHeaderKey, true))

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

func (c *Converter) tableCaption(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

const (
	// Marker to indicate we are inside the header
	insideTableHeaderKey = "inside.table.header"
)

func insideTableHeader(ctx context.Context) bool {
	return ctx.Value(insideTableHeaderKey) != nil
}

func (c *Converter) tr(n *html.Node, ctx context.Context) error {

	// The table definition
	table := stylesheet.TableFromContext(ctx)

	cell := 0
	firstCell := true
	err := parser.HandleChildren(func(n *html.Node, ctx context.Context) (err error) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "th", "td":

				cellDef := table.GetColumn(cell)
				if cellDef.Hidden {
					cell++
					return
				}

				// Handle cell separator
				if firstCell {
					firstCell = false
				} else {
					err = WriteString(ctx, " & ")
				}

				style := c.Stylesheet().GetStyle(n)

				colSpan, _ := parser.GetAttrInt(n, "colspan", 1)
				width := stylesheet.FloatToString(
					table.GetColumnWidth(cell, colSpan),
					table.GetColumn(cell).ColUnit)
				rowSpan, _ := parser.GetAttrInt(n, "rowspan", 1)
				multiCol, multiRow := colSpan > 1, rowSpan > 1

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
						// Ignore whitespance only text nodes
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
						stylesheet.DefString(width, "l"),
					)
				}
				if err == nil && multiRow {
					err = Writef(ctx, `\multirow{%d}{%s}{`,
						rowSpan,
						// width or "*" for contents natural width
						stylesheet.DefString(width, "*"),
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
				// FIXME this might not account for intermediate hidden columns
				cell = cell + colSpan
			default:
			}
		}
		return
	}, n, ctx)

	//if err == nil && insideTableHeader(ctx) {
	//	// This at the end of each header row, so we repeat headers when the
	//	// table spans multiple pages
	//	err = WriteStringLn(ctx, `\endhead`)
	//}

	if err == nil {
		err = WriteStringLn(ctx, ` \\`)
	}

	if err == nil && insideTableHeader(ctx) {
		err = WriteStringLn(ctx, `\endhead`)
	}

	return err
}
