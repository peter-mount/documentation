package latex

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"io"
	"unsafe"
)

func New() parser.Handler {

	// Handle normal html content
	content := parser.New()

	// This handler is used by multiple elements
	heading := parser.Of(headingStart, content.Handler().HandleChildren, headingEnd)

	content.
		Text(text).
		// anchor is ignored, just parse it's content
		Handle("a", handleChildren).
		Handle("br", lineBreak).
		Handle("div", div).
		Handle("p", paragraph).
		Handle("figure", parser.Of(
			figureStart,
			parser.New().
				Handle("figcaption", parser.Of(figureCaptionStart, content.Handler().HandleChildren, figureCaptionEnd)).
				Default(content.Handler()).
				Handler().
				HandleChildren,
			figureEnd1,
			//contentHandler.Type("figcaption").HandleChildren,
			figureEnd2)).
		Handle("h1", heading).
		Handle("h2", heading).
		Handle("h3", heading).
		Handle("h4", heading).
		Handle("h5", heading).
		Handle("table", parser.Of(
			tableStart,
			parser.New().
				Handle("thead", tableHead).
				Handle("tbody", handleChildren).
				Handle("tr", tr).
				// FIXME Text & Default need to be within tr only
				// but until I figure a way to handle "&" separators with Scanner this will have to do
				Text(text).
				Default(content.Handler()).
				Handler().
				HandleChildren,
			tableEnd1,
			parser.Of(tableCaption).Type("caption").HandleChildren,
			tableEnd2)).
		// These elements are ignored
		Handle("bookMeta", nil)

	// Final Handler
	return parser.Of(
		beginDocument,
		parser.New().
			Handle("div",
				content.Handler().
					HasClasses("td-content", "td-content-new").
					HandleChildren).
			Default(handleChildren).
			Handler().
			FindByClass("td-main", "td-main-new"),
		endDocument)
}

func WithContext(w io.Writer, ctx context.Context) context.Context {
	return context.WithValue(ctx, "writer", w)
}

func Writer(ctx context.Context) io.Writer {
	return ctx.Value("writer").(io.Writer)
}

func Write(ctx context.Context, b ...byte) error {
	_, err := Writer(ctx).Write(b)
	return err
}

func WriteString(ctx context.Context, s string) error {
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	return Write(ctx, b...)
}

func Writef(ctx context.Context, f string, a ...any) error {
	return WriteString(ctx, fmt.Sprintf(f, a...))
}

func WriteStringLn(ctx context.Context, s string) error {
	return WriteString(ctx, s+"\n")
}

func Comment(ctx context.Context, f string, a ...any) error {
	return Writef(ctx, "%% "+f+"\n", a...)
}

func handleChildren(n *html.Node, ctx context.Context) error {
	return parser.HandleChildren(parser.ScannerFromContext(ctx).Handle, n, ctx)
}

func handleBracedChildren(n *html.Node, ctx context.Context) error {
	return handleWrapped('{', '}', n, ctx)
}

func handleWrapped(s, e byte, n *html.Node, ctx context.Context) error {
	err := Write(ctx, s)

	if err == nil {
		err = handleChildren(n, ctx)
	}

	if err == nil {
		err = Write(ctx, e)
	}

	return err
}
