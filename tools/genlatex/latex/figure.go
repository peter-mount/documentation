package latex

import (
	"context"
	"golang.org/x/net/html"
)

func figureStart(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\begin{figureOuter}\n\\begin{figureInner}\n")
}

func figureEnd1(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\end{figureInner}\n")
}

func figureEnd2(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\end{figureOuter}\n")
}

func figureCaptionStart(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	//err := WriteString(ctx, "\\figure-caption{\n")
	//if err == nil {
	err := handleChildren(n, ctx)
	//}
	return err
}

func figureCaptionEnd(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "}\n")
}
