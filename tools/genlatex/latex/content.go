package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"golang.org/x/net/html"
)

func (c *Converter) text(n *html.Node, ctx context.Context) error {
	s := n.Data
	//s := strings.TrimSpace(n.Data)
	if s != "" {
		return util.WriteString(ctx, EscapeText(s))
	}
	return nil
}

func (c *Converter) paragraph(n *html.Node, ctx context.Context) error {
	// Paragraphs can also be marginNotes
	if parser.HasClass(n, "marginNote") {
		return c.marginNote(n, ctx)
	}
	if parser.HasClass(n, "sideNote") {
		return c.sideNote(n, ctx)
	}

	err := util.HandleChildren(n, ctx)
	if err == nil {
		err = util.WriteStringLn(ctx, "\n")
	}
	return err
}

// Appends \\ to the end of the current line to indicate a line break
func (c *Converter) lineBreak(_ *html.Node, ctx context.Context) error {
	if stylesheet.TableFromContext(ctx) != nil {
		return util.WriteStringLn(ctx, `\\`)
	}
	return util.WriteStringLn(ctx, `\hfill \break`)
}

// code handles the em html element
func (c *Converter) em(n *html.Node, ctx context.Context) error {
	return util.HandleSimpleCommand(`\textit`, n, ctx)
}

// code handles the em html element
func (c *Converter) strong(n *html.Node, ctx context.Context) error {
	return util.HandleSimpleCommandSpace(`\textbf`, n, ctx)
}

// code handles the em html element
func (c *Converter) sup(n *html.Node, ctx context.Context) error {
	return util.HandleSimpleCommand(`\textsuperscript`, n, ctx)
}
