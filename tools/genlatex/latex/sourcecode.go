package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

// sourceCode handles code listings
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

// code handles the code html element
func (c *Converter) code(n *html.Node, ctx context.Context) error {
	err := Write(ctx, '~')
	if err == nil {
		err = handleSimpleCommand(`\texttt`, n, ctx)
	}
	if err == nil {
		err = Write(ctx, '~')
	}
	return err
}
