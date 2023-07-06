package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

func div(n *html.Node, ctx context.Context) error {
	_ = Comment(ctx, "div class=%q", parser.GetAttr(n, "class"))
	return handleChildren(n, ctx)
}

func text(n *html.Node, ctx context.Context) error {
	s := strings.TrimSpace(n.Data)
	if s != "" {
		s = strings.ReplaceAll(s, "&", "\\&")

		return WriteString(ctx, s)
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
	return WriteStringLn(ctx, " \\linebreak")
}
