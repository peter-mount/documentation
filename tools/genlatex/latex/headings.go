package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

// heading handles headings h1 h2 etc. based on class type
func (c *Converter) headingStart(n *html.Node, ctx context.Context) error {
	var err error

	hType := n.Data

	switch {
	// title is for a non-numbered heading?
	case hType == "h1" && parser.HasClass(n, "title"):
		err = WriteString(ctx, `\chapter`)

	case hType == "h1":
		err = WriteString(ctx, `\section`)

	case hType == "h2":
		err = WriteString(ctx, `\subsection`)

	case hType == "h3":
		err = WriteString(ctx, `\subsubsection`)

	case hType == "h4",
		hType == "h5":
		err = WriteString(ctx, `\paragraph`)

		// Default numbered heading
	default:
		err = WriteString(ctx, `\section`)
	}

	if err == nil {
		err = Write(ctx, '{')
	}

	return err
}

func (c *Converter) headingEnd(n *html.Node, ctx context.Context) error {
	return Write(ctx, '}', '\n')
}
