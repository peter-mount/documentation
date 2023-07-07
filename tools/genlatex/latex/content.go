package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

func div(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func text(n *html.Node, ctx context.Context) error {
	s := strings.TrimSpace(n.Data)
	if s != "" {
		return WriteString(ctx, EscapeText(s))
	}
	return nil
}

func paragraph(n *html.Node, ctx context.Context) error {
	err := handleChildren(n, ctx)
	if err == nil {
		err = WriteStringLn(ctx, "\n")
	}
	return err
}

// Appends \\ to the end of the current line to indicate a line break
func lineBreak(_ *html.Node, ctx context.Context) error {
	if insideTable(ctx) {
		return WriteStringLn(ctx, `\\`)
	}
	return WriteStringLn(ctx, ` \linebreak`)
}

func ul(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "print-page-link") {
		return WriteStringLn(ctx, `\toc`)
	}

	return handleChildren(n, ctx)
}
