package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

// heading handles headings h1 h2 etc. based on class type
func headingStart(n *html.Node, ctx context.Context) error {
	var err error

	hType := n.Data

	switch {
	// title is for a non-numbered heading?
	case hType == "h1" && parser.HasClass(n, "title"):
		err = WriteString(ctx, `\h*`)

	case hType == "h1":
		err = WriteString(ctx, `\h`)

	case hType == "h2":
		err = WriteString(ctx, `\hh`)

	case hType == "h3":
		err = WriteString(ctx, `\hhh`)

	case hType == "h4":
		err = WriteString(ctx, `\hhh`)

	case hType == "h5":
		err = WriteString(ctx, `\hhh`)

		// Default numbered heading
	default:
		err = WriteString(ctx, `\h`)
	}

	if err == nil {
		err = Write(ctx, '{')
	}

	return err
}

func headingEnd(n *html.Node, ctx context.Context) error {
	return Write(ctx, '}', '\n')
}
