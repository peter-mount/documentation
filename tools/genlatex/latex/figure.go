package latex

import (
	"context"
	"golang.org/x/net/html"
)

func (c *Converter) figureStart(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\begin{figureOuter}\n\\begin{figureInner}\n")
}

func (c *Converter) figureEnd1(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\end{figureInner}\n")
}

func (c *Converter) figureEnd2(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "\\end{figureOuter}\n")
}

func (c *Converter) figureCaptionStart(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	//err := WriteString(ctx, "\\figure-caption{\n")
	//if err == nil {
	//err := handleChildren(n, ctx)
	//}
	//return err
	return nil
}

func (c *Converter) figureCaptionEnd(n *html.Node, ctx context.Context) error {
	// FIXME these are temp holders until I find out how LaTeX handles figures
	return nil //WriteString(ctx, "}\n")
}
