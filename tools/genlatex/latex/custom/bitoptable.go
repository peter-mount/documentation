package custom

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func BitOpTable(n *html.Node, ctx context.Context) error {
	if err := util.Writef(ctx, `\bitOpTable{%s}{`, parser.GetTextByClass(n, "bitOpTableCaption", "")); err != nil {
		return err
	}

	// The header row
	if c := parser.FindByClass(n, "bitOpTableHeader"); c != nil {
		err := util.WriteString(ctx, `\textbf{`)
		if err == nil {
			err = bitOpTableRow(c, ctx, `\bitOpTableCell`, `\bitOpTableCellA`, `\bitOpTableCellB`)
		}
		if err == nil {
			err = util.Write(ctx, '}')
		}

		if err != nil {
			return err
		}
	}

	// The body rows, just find the first one then run against each sibling
	for c := parser.FindByClass(n, "bitOpTableRow"); c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tr" {
			if err := bitOpTableRow(c, ctx, `\bitOpTableCell`, `\bitOpTableCellC`, `\bitOpTableCellD`); err != nil {
				return err
			}
		}
	}

	return util.Write(ctx, '}')
}

func bitOpTableRow(n *html.Node, ctx context.Context, c1, c2, c3 string) error {
	// Slice of TD elements
	var td []*html.Node
	c := parser.FindElement(n, "td")
	if c == nil {
		c = parser.FindElement(n, "th")
	}
	for ; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			td = append(td, c)
		}
	}

	tdl := len(td)
	if tdl < 5 {
		return fmt.Errorf("too short a bitOpTable row, got %d expected >4", tdl)
	}

	err := util.WriteString(ctx, `\bitOpTableRow`)
	// First two direct in command args
	for i := 0; i < 2; i++ {
		if err == nil {
			err = util.HandleBracedChildren(td[i], ctx)
		}
	}

	if err == nil {
		err = util.WriteString(ctx, `{`)
	}

	if err == nil {
		for i := 2; i < tdl; i++ {
			switch i {
			case 2, tdl - 1:
				err = util.WriteString(ctx, c1)

			case tdl - 2:
				err = util.WriteString(ctx, c3)

			default:
				err = util.WriteString(ctx, c2)
			}
			if err == nil {
				err = util.HandleBracedChildren(td[i], ctx)
			}
		}
	}

	if err == nil {
		err = util.WriteString(ctx, `}`)
	}

	return err
}
