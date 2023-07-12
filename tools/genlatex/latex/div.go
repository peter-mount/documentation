package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func (c *Converter) div(n *html.Node, ctx context.Context) error {
	switch {
	case parser.HasClass(n, "lead"):
		return environment("huge", n, ctx)

	case parser.HasClass(n, "sourceCode"):
		return c.sourceCode(n, ctx)

	case parser.HasClass(n, "marginNote"):
		return c.marginNote(n, ctx)

	default:
		return handleChildren(n, ctx)
	}
}

func (c *Converter) marginNote(n *html.Node, ctx context.Context) error {
	// \marginnote[offset]{text}
	err := WriteString(ctx, `\marginnote`)
	// TODO add optional [offset[ as an optional attribute
	if err == nil {
		err = handleBracedChildren(n, ctx)
	}
	return nil
}
