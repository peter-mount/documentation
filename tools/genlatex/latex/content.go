package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"golang.org/x/net/html"
	"strings"
)

func (c *Converter) text(n *html.Node, ctx context.Context) error {
	s := strings.TrimSpace(n.Data)
	if s != "" {
		return WriteString(ctx, EscapeText(s))
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

	err := handleChildren(n, ctx)
	if err == nil {
		err = WriteStringLn(ctx, "\n")
	}
	return err
}

// Appends \\ to the end of the current line to indicate a line break
func (c *Converter) lineBreak(_ *html.Node, ctx context.Context) error {
	if stylesheet.TableFromContext(ctx) != nil {
		return WriteStringLn(ctx, `\\`)
	}
	return WriteStringLn(ctx, `\hfill \break`)
}

// code handles the em html element
func (c *Converter) em(n *html.Node, ctx context.Context) error {
	return handleSimpleCommand(`\textsl`, n, ctx)
}

// code handles the em html element
func (c *Converter) strong(n *html.Node, ctx context.Context) error {
	return handleSimpleCommand(`\textbf`, n, ctx)
}

// code handles the em html element
func (c *Converter) sup(n *html.Node, ctx context.Context) error {
	return handleSimpleCommand(`\textsuperscript`, n, ctx)
}
