package latex

import (
	"context"
	"golang.org/x/net/html"
)

func tableStart(n *html.Node, ctx context.Context) error {
	return WriteString(ctx, "\\begin{table}[h]\n\\begin{tabular}{ |l|l| }\n")
}

func tableEnd1(n *html.Node, ctx context.Context) error {
	return WriteString(ctx, "\\end{tabular}\n")
}

func tableEnd2(n *html.Node, ctx context.Context) error {
	return WriteString(ctx, "end{table}\n")
}

func tableCaption(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func thead(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func tbody(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func tr(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func thStart(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func thEnd(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func tdStart(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}

func tdEnd(n *html.Node, ctx context.Context) error {
	return handleChildren(n, ctx)
}
