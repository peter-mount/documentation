package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func (c *Converter) ul(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "print-page-link") {
		return util.WriteStringLn(ctx, `\toc`)
	}

	return util.Environment("itemize", n, ctx)
}

func (c *Converter) ol(n *html.Node, ctx context.Context) error {
	return util.Environment("enumerate", n, ctx)
}

func (c *Converter) li(n *html.Node, ctx context.Context) error {
	err := util.WriteString(ctx, "\n\\item{}")
	if err == nil {
		err = util.HandleChildren(n, ctx)
	}
	return err
}
