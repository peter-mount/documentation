package util

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func getString(n string, ctx context.Context) string {
	v := ctx.Value(n)
	if v == nil {
		return ""
	}
	return v.(string)
}

func WithCell(ctx context.Context, width, height string) context.Context {
	return WithCellWidth(WithCellHeight(ctx, height), width)
}

func WithCellWidth(ctx context.Context, width string) context.Context {
	return context.WithValue(ctx, "cellWidth", width)
}

func CellWidth(ctx context.Context) string {
	return getString("cellWidth", ctx)
}

func WithCellHeight(ctx context.Context, height string) context.Context {
	return context.WithValue(ctx, "cellHeight", height)
}

func CellHeight(ctx context.Context) string {
	return getString("cellHeight", ctx)
}

func CellRow(h parser.Handler, n *html.Node, ctx context.Context) error {

	if err := WriteString(ctx, `\hexmaprow{}{`); err != nil {
		return err
	}

	if err := parser.HandleChildren(h, n, ctx); err != nil {
		return err
	}

	var a []string
	a = append(a, "}")

	return WriteSlice(ctx, a)
}

func Cell(borders, colour string, h parser.Handler, n *html.Node, ctx context.Context) error {
	return CFBox(borders, colour, func(n *html.Node, ctx context.Context) error {
		err := Writef(ctx, `\vbox to %s {\hbox to %s {`, CellHeight(ctx), CellWidth(ctx))

		if err == nil {
			err = h.Do(n, ctx)
		}

		if err == nil {
			err = WriteString(ctx, "}}")
		}

		return nil
	}, n, ctx)
}
