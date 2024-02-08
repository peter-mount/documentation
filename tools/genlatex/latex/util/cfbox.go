package util

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

// CFBox allows us to handle cell borders
//
// cfbox which allows a coloured border around a box to be drawn based on which border
// we want see https://tex.stackexchange.com/a/55534
func CFBox(borders, colour string, h parser.Handler, n *html.Node, ctx context.Context) error {
	if colour == "" {
		colour = "black"
	}

	err := Writef(ctx, `\cfbox[%s]{%s}{`, borders, colour)

	if err == nil {
		err = h.Do(n, ctx)
	}

	if err == nil {
		err = WriteString(ctx, "}")
	}

	return err
}
