package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func (c *Converter) ul(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "print-page-link") {
		return WriteStringLn(ctx, `\toc`)
	}

	return environment("itemize", n, ctx)
}

func (c *Converter) ol(n *html.Node, ctx context.Context) error {
	return environment("enumerate", n, ctx)
}

func (c *Converter) li(n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\item{}")
	if err == nil {
		err = handleChildren(n, ctx)
	}
	return err
}
