package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

func (c *Converter) div(n *html.Node, ctx context.Context) error {
	switch {
	case parser.HasClass(n, "lead"):
		return environment("huge", n, ctx)

	case parser.HasClass(n, "sourceCode"):
		return c.sourceCode(n, ctx)

	default:
		return handleChildren(n, ctx)
	}
}

func (c *Converter) text(n *html.Node, ctx context.Context) error {
	s := strings.TrimSpace(n.Data)
	if s != "" {
		return WriteString(ctx, EscapeText(s))
	}
	return nil
}

func (c *Converter) paragraph(n *html.Node, ctx context.Context) error {
	err := handleChildren(n, ctx)
	if err == nil {
		err = WriteStringLn(ctx, "\n")
	}
	return err
}

// Appends \\ to the end of the current line to indicate a line break
func (c *Converter) lineBreak(_ *html.Node, ctx context.Context) error {
	if insideTable(ctx) {
		return WriteStringLn(ctx, `\\`)
	}
	return WriteStringLn(ctx, `\hfill \break`)
}
