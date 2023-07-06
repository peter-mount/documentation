package latex

import (
	"context"
	"golang.org/x/net/html"
)

func head(n *html.Node, ctx context.Context) error {
	return nil
}

func beginDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `\begin{document}`)
}

func endDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `\end{document}`)
}
