package util

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"io"
	"unsafe"
)

func Comment(ctx context.Context, f string, a ...any) error {
	return Writef(ctx, "%% "+f+"\n", a...)
}

func Group(f parser.Handler, n *html.Node, ctx context.Context) error {
	err := WriteString(ctx, "\n\\begingroup\n")
	if err == nil {
		err = f(n, ctx)
	}
	if err == nil {
		err = WriteString(ctx, "\n\\endgroup\n")
	}
	return err
}

func Environment(t string, n *html.Node, ctx context.Context) error {
	err := Writef(ctx, "\n\\begin{%s}\n", t)
	if err == nil {
		err = HandleChildren(n, ctx)
	}
	if err == nil {
		err = Writef(ctx, "\n\\end{%s}\n", t)
	}
	return err
}

func WithContext(w io.Writer, ctx context.Context) context.Context {
	return context.WithValue(ctx, "writer", w)
}

func Writer(ctx context.Context) io.Writer {
	return ctx.Value("writer").(io.Writer)
}

func Write(ctx context.Context, b ...byte) error {
	if len(b) == 0 {
		return nil
	}

	_, err := Writer(ctx).Write(b)
	return err
}

func WriteString(ctx context.Context, s string) error {
	if s == "" {
		return nil
	}

	b := unsafe.Slice(unsafe.StringData(s), len(s))
	return Write(ctx, b...)
}

func Writef(ctx context.Context, f string, a ...any) error {
	return WriteString(ctx, fmt.Sprintf(f, a...))
}

func WriteStringLn(ctx context.Context, s string) error {
	return WriteString(ctx, s+"\n")
}

func WriteSlice(ctx context.Context, s []string) error {
	for _, a := range s {
		if err := WriteStringLn(ctx, a); err != nil {
			return err
		}
	}
	return nil
}
