package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func head(n *html.Node, ctx context.Context) error {
	return nil
}

func beginDocument(n *html.Node, ctx context.Context) error {
	_ = WriteStringLn(ctx,
		`\documentclass{textbook}
\usepackage{multirow}
\usepackage{array}
\usepackage{longtable}
\lang      {english}
`)

	// Look for bookMeta "object"
	meta := parser.FindById(n, "bookMeta")
	if meta != nil {
		_ = parser.HandleChildren(func(n *html.Node, ctx context.Context) error {
			if n.Type == html.ElementNode {
				key := n.Data
				value := parser.GetText(n)
				switch key {
				case "cover":
					key = key + "*"
					// don't escape the text
					value = "{" + value + "}"

				case "dedicate", "edition", "license":
					// Don't escape the text or change key

				default:
					value = "{" + escapeText(value) + "}"
				}
				return Writef(ctx, "\\%s %s\n", key, value)
			}
			return nil
		}, meta, ctx)
	}

	return WriteStringLn(ctx, `\begin{document}`)
}

func endDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `\end{document}`)
}
