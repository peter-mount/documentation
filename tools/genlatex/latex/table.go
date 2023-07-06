package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

func tableStart(n *html.Node, ctx context.Context) error {
	// We need to calculate the number of columns in the table
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

	return Writef(ctx, "\n\\begin{tabular}{%s}\n", strings.Repeat("l", cols))
}

func tableEnd1(n *html.Node, ctx context.Context) error {
	return WriteString(ctx, "\\end{tabular}\n")
}

func tableEnd2(n *html.Node, ctx context.Context) error {
	return nil //WriteString(ctx, "end{table}\n")
}

func tableCaption(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func tr(n *html.Node, ctx context.Context) error {
	cell := 1
	err := parser.HandleChildren(func(n *html.Node, ctx context.Context) (err error) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "th", "td":
				if cell > 1 {
					err = WriteString(ctx, " & ")
				}

				align := getAlign(n)
				colSpan, _ := parser.GetAttrInt(n, "colspan", 1)
				rowSpan, _ := parser.GetAttrInt(n, "rowspan", 1)
				multiCol, multiRow := colSpan > 1 || align != "", rowSpan > 1

				if err == nil && multiCol {
					err = Writef(ctx, `\multicolumn{%d}{%s}{`, colSpan, align)
				}

				if err == nil && multiRow {
					err = Writef(ctx, `\multirow{%d}{*}{`, rowSpan)
				}

				if err == nil {
					err = handleChildren(n, ctx)
				}
				if err == nil && multiRow {
					err = Write(ctx, '}')
				}
				if err == nil && multiCol {
					err = Write(ctx, '}')
				}

				cell++
			default:
			}
		}
		return
	}, n, ctx)
	if err == nil {
		err = WriteStringLn(ctx, `\\`)
	}
	return err
}

func getAlign(n *html.Node) string {
	for _, class := range parser.GetClass(n) {
		switch class {
		case "offset":
			return "r"
		}
	}
	return ""
}
