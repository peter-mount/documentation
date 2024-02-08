package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func (c *Converter) div(n *html.Node, ctx context.Context) error {
	switch {
	case parser.HasClass(n, "printPageBreakAvoid"):
		return util.Environment("samepage", n, ctx)

	case parser.HasClass(n, "lead"):
		return util.Environment("huge", n, ctx)

	case parser.HasClass(n, "sourceCode"):
		return c.sourceCode(n, ctx)

	case parser.HasClass(n, "marginNote"):
		return c.marginNote(n, ctx)

	case parser.HasClass(n, "sideNote"):
		return c.sideNote(n, ctx)

	default:
		return util.HandleChildren(n, ctx)
	}
}

func (c *Converter) marginNote(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "tableAlign") {
		return util.BuffersFromContext(ctx).
			UseBuffer("marginNote", c.marginNoteImpl, n, ctx)
	}
	return c.marginNoteImpl(n, ctx)
}

func (c *Converter) marginNoteImpl(n *html.Node, ctx context.Context) error {
	// \marginnote[offset]{text}
	err := util.WriteString(ctx, `\marginnote`)
	// TODO add optional [offset[ as an optional attribute
	if err == nil {
		err = util.HandleBracedChildren(n, ctx)
	}
	return nil
}

func (c *Converter) sideNote(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "tableAlign") {
		return util.BuffersFromContext(ctx).
			UseBuffer("sideNote", c.sideNoteImpl, n, ctx)
	}
	return c.sideNoteImpl(n, ctx)
}

func (c *Converter) sideNoteImpl(n *html.Node, ctx context.Context) error {
	// \sidenote[offset]{text}
	err := util.WriteString(ctx, `\sidenote`)
	// TODO add optional [offset[ as an optional attribute
	if err == nil {
		err = util.HandleBracedChildren(n, ctx)
	}
	return nil
}
