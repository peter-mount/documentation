package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func (c *Converter) sourceCode(n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\begin{lstlisting}\n")

	if err == nil {
		err = parser.HandleChildren(func(n *html.Node, ctx context.Context) error {
			if n.Type == html.TextNode {
				return WriteStringLn(ctx, n.Data)
			}
			return nil
		}, n, ctx)
	}

	if err == nil {
		err = WriteString(ctx, "\n\\end{lstlisting}\n")
	}
	return err
}
