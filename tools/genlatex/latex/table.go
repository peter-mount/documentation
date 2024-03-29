package latex

import (
	"bytes"
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/custom"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
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
		_ = util.Write(ctx, ' ', '\\', '\\', '\n')
	}

	switch {
	case parser.HasClass(n, "bitOpTable"):
		return custom.BitOpTable(n, ctx)

	case parser.HasClass(n, "m6502opcode"):
		return custom.OpcodeTable6502(n, ctx)

	case parser.HasClass(n, "processorFlags"):
		return custom.ProcessorFlags(n, ctx)

	case parser.HasClass(n, "truthTable"):
		return custom.TruthTable(n, ctx)

	default:
		// Needed to save state otherwise headers spanning pages will break subsequent tables
		return util.Group(c.tableImpl, n, ctx)
	}
}

func (c *Converter) tableImpl(n *html.Node, ctx context.Context) error {

	// if nested then tabular otherwise the outer table is a longtable
	tType := "longtblr"
	if stylesheet.TableFromContext(ctx) != nil {
		tType = "tblr"

		// Reset header flag incase we are also in the thead section of the parent
		ctx = context.WithValue(ctx, insideTableHeaderKey, nil)
	}

	// how many columns are we?
	cols := getTableColCount(n, ctx)

	//style := c.Stylesheet().GetStyle(n)
	table := c.Stylesheet().GetTable(n)
	ctx = table.WithContext(ctx)

	if table.FontSize != "" {
		if err := util.WriteString(ctx, table.FontSize); err != nil {
			return err
		}
	}

	// Initialise table state
	ctx = newTableState(tType, cols, table, ctx)
	model := tableStateFromContext(ctx).table

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

			case "caption":
				caption, err := util.HandleChildrenString(n, ctx)
				if err != nil {
					return err
				}
				model.Caption = caption
			}
		}
		return f.HandleChildren(n, ctx)
	}
	err := f.HandleChildren(n, ctx)

	if err == nil {
		model.Finalise()
		err = util.WriteString(ctx, model.String())
	}

	return err
}

func (c *Converter) tableHeader(f parser.Handler, n *html.Node, ctx context.Context) error {
	// Mark we are in a header
	ctx = context.WithValue(ctx, insideTableHeaderKey, true)

	//table := stylesheet.TableFromContext(ctx)
	var err error

	//if table.HeaderPrefix != "" {
	//	err = WriteStringLn(ctx, table.HeaderPrefix)
	//}

	if err == nil {
		err = parser.HandleChildren(f, n, ctx)
	}

	//if err == nil && table.HeaderSuffix != "" {
	//	err = WriteStringLn(ctx, table.HeaderSuffix)
	//}

	if err == nil {
		//err = WriteStringLn(ctx, `\endhead`)
	}

	return err
}

func (c *Converter) tableCaption(n *html.Node, ctx context.Context) error {
	return nil // TODO HandleChildren(n, ctx)
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
		return false, util.WriteString(ctx, " & ")
	}
}

func (c *Converter) tr(n *html.Node, ctx context.Context) error {

	// The table definition
	//table := stylesheet.TableFromContext(ctx)
	ts := tableStateFromContext(ctx)

	// append to the last row in the model
	row := ts.table.RowCount()
	col := 0

	err := parser.HandleChildren(func(n *html.Node, ctx context.Context) (err error) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "th", "td":

				// Skip columns that are row-spanning
				for ts.decRowSpan(col) {
					col++
				}

				//style := c.Stylesheet().GetStyle(n)
				colSpan, _ := parser.GetAttrInt(n, "colspan", 1)
				rowSpan, _ := parser.GetAttrInt(n, "rowspan", 1)

				// Look ahead, if a table exists then we need to wrap the col to allow line break before it.
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

				tableCell := &TableCell{
					ColSpan: colSpan,
					RowSpan: rowSpan,
				}
				ts.table.SetCell(row, col, tableCell)

				var buf bytes.Buffer
				cellCtx := util.WithContext(&buf, ctx)

				if err == nil && tabularCell {
					// Wrap col content within a tabular block, so we can then use \\ as a line break
					err = util.Writef(cellCtx,
						`\begin{tblr}[%s]{@{}%s@{}}`,
						"t",
						stylesheet.DefStrings(
							ts.table.Table.GetColumnDef(col, colSpan),
							"l",
						),
					)
				}

				if err == nil {
					err = util.HandleChildren(n, cellCtx)
				}

				if err == nil && tabularCell {
					// Wrap col content within a tabular block, so we can then use \\ as a line break
					err = util.WriteString(cellCtx, `\end{tblr}`)
				}

				tableCell.Text = buf.String()

				// Increment col number
				for i := 0; i < colSpan; i++ {
					ts.setRowSpan(col, rowSpan)
					col++
				}
			default:
			}
		}
		return
	}, n, ctx)

	// Clean up the row
	ts.finishRow(col)

	return err
}
