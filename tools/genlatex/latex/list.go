package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func ul(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "print-page-link") {
		return WriteStringLn(ctx, `\toc`)
	}

	return environment("itemize", n, ctx)
}

func ol(n *html.Node, ctx context.Context) error {
	return environment("enumerate", n, ctx)
}

func li(n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\item{}")
	if err == nil {
		err = handleChildren(n, ctx)
	}
	return err
}
