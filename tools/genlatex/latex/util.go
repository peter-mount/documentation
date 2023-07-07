package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func group(f parser.Handler, n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\begingroup{}\n")
	if err == nil {
		err = f(n, ctx)
	}
	if err == nil {
		err = WriteString(ctx, "\n\\endgroup{}\n")
	}
	return err
}

func environment(t string, n *html.Node, ctx context.Context) error {
	err := Writef(ctx, "\n\\begin{%s}\n", t)
	if err == nil {
		err = handleChildren(n, ctx)
	}
	if err == nil {
		err = Writef(ctx, "\n\\end{%s}\n", t)
	}
	return err
}
