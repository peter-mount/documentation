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

	return list("itemize", n, ctx)
}

func ol(n *html.Node, ctx context.Context) error {
	return list("enumerate", n, ctx)
}

func list(t string, n *html.Node, ctx context.Context) error {
	err := Writef(ctx, "\n\\begingroup{}\\begin{%s}", t)
	if err == nil {
		err = handleChildren(n, ctx)
	}
	if err == nil {
		err = Writef(ctx, "\n\\end{%s}\\endgroup{}\n", t)
	}
	return err
}

func li(n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\item{}")
	if err == nil {
		err = handleChildren(n, ctx)
	}
	return err
}
