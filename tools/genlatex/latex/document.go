package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"time"
)

func (c *Converter) beginDocument(n *html.Node, ctx context.Context) error {
	s := c.Stylesheet()

	_ = Writef(ctx, "%% Generated %s\n", time.Now().Format(time.RFC3339))

	_ = Writef(ctx, "\\documentclass{%s}\n", s.DocumentClass)
	for _, p := range s.UsePackage {
		_ = Writef(ctx, "\\usepackage{%s}\n", p)
	}

	_ = WriteString(ctx, `
\usepackage{array}
\newcolumntype{H}{>{\setbox0=\hbox\bgroup}c<{\egroup}@{}}
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

				// These append "*" and wrap with () but do not escape
				case "cover":
					key = key + "*"
					value = "{" + value + "}"

				// Don't escape the text or change key
				case "dedicate", "edition", "license":

				// Wrap with {} but do not escape and do not change key
				case "isbn":
					value = "{" + value + "}"

				// Default wrap {} and escape the value
				default:
					value = "{" + EscapeText(value) + "}"
				}
				return Writef(ctx, "\\%s %s\n", key, value)
			}
			return nil
		}, meta, ctx)
	}

	return WriteStringLn(ctx, `

\lstset{
	basicstyle=\footnotesize\ttfamily,
	xleftmargin=\parindent,
	frame=L,
	breaklines=true,
}

\begin{document}
`)
}

func (c *Converter) endDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `
\end{document}
`)
}
