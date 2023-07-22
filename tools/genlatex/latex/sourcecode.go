package latex

import (
	"bytes"
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

// sourceCode handles code listings
func (c *Converter) sourceCode(n *html.Node, ctx context.Context) error {
	err := util.WriteString(ctx, "\n\\begin{lstlisting}\n")

	var buf bytes.Buffer
	srcCtx := util.WithContext(&buf, ctx)

	if err == nil {
		err = parser.HandleChildren(func(n *html.Node, ctx context.Context) error {
			if n.Type == html.TextNode {
				return util.WriteStringLn(ctx, n.Data)
			}
			return nil
		}, n, srcCtx)
	}

	if err == nil {
		err = util.WriteString(ctx, strings.Trim(buf.String(), "\n\r"))
	}

	if err == nil {
		err = util.WriteString(ctx, "\n\\end{lstlisting}\n")
	}

	return err
}

// code handles the code html element
func (c *Converter) code(n *html.Node, ctx context.Context) error {
	err := util.Write(ctx, '~')
	if err == nil {
		err = util.HandleSimpleCommand(`\texttt`, n, ctx)
	}
	if err == nil {
		err = util.Write(ctx, '~')
	}
	return err
}
