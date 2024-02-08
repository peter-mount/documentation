package util

import (
	"bytes"
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func HandleChildren(n *html.Node, ctx context.Context) error {
	return parser.HandleChildren(parser.ScannerFromContext(ctx).Handle, n, ctx)
}

func HandleChildrenAsString(n *html.Node, ctx context.Context) (string, error) {
	var buf bytes.Buffer
	cellCtx := WithContext(&buf, ctx)
	if err := HandleChildren(n, cellCtx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// HandleChildrenString is the same as HandleChildren except it returns the
// content as a string rather than write it to the stream.
func HandleChildrenString(n *html.Node, ctx context.Context) (string, error) {
	var buf bytes.Buffer
	err := HandleChildren(n, WithContext(&buf, ctx))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func HandleBracedChildren(n *html.Node, ctx context.Context) error {
	return HandleWrapped('{', '}', n, ctx)
}

func HandleWrapped(s, e byte, n *html.Node, ctx context.Context) error {
	err := Write(ctx, s)

	if err == nil {
		err = HandleChildren(n, ctx)
	}

	if err == nil {
		err = Write(ctx, e)
	}

	return err
}

func HandleSimpleCommand(c string, n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, c)
	if err == nil {
		err = HandleBracedChildren(n, ctx)
	}
	return err
}

// HandleSimpleCommandSpace is the same as HandleSimpleCommand but appends a space to the output
func HandleSimpleCommandSpace(c string, n *html.Node, ctx context.Context) error {
	err := HandleSimpleCommand(c, n, ctx)
	if err == nil {
		err = Write(ctx, ' ')
	}
	return err
}
