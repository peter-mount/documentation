package latex

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"unsafe"
)

type Converter struct {
	handler    parser.Handler
	stylesheet *stylesheet.Stylesheet
}

func (c *Converter) Handler() parser.Handler {
	return c.handler
}

func (c *Converter) Stylesheet() *stylesheet.Stylesheet {
	return c.stylesheet
}

func New(config string) (*Converter, error) {
	c := &Converter{}

	// load stylesheet
	b, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	c.stylesheet = &stylesheet.Stylesheet{}
	err = yaml.Unmarshal(b, &c.stylesheet)
	if err != nil {
		return nil, err
	}

	// Ensure it's fully configured
	c.stylesheet.Init()

	// Handle normal html content
	content := parser.New()

	// This handler is used by multiple elements
	heading := parser.Of(c.headingStart, content.Handler().HandleChildren, c.headingEnd)

	content.
		Text(c.text).
		// anchor is ignored, just parse it's content
		Handle("a", handleChildren).
		Handle("br", c.lineBreak).
		Handle("code", handleChildren).
		Handle("div", c.div).
		Handle("dd", handleChildren).
		Handle("dl", handleChildren).
		Handle("dt", handleChildren).
		Handle("em", handleChildren).
		Handle("figure", parser.Of(
			c.figureStart,
			parser.New().
				Handle("figcaption", nil).
				//Handle("figcaption", parser.Of(c.figureCaptionStart, content.Handler().HandleChildren, c.figureCaptionEnd)).
				Default(content.Handler()).
				Handler().
				HandleChildren,
			c.figureEnd1,
			//contentHandler.Type("figcaption").HandleChildren,
			c.figureEnd2)).
		Handle("h1", heading).
		Handle("h2", heading).
		Handle("h3", heading).
		Handle("h4", heading).
		Handle("h5", heading).
		Handle("li", c.li).
		Handle("ol", c.ol).
		Handle("p", c.paragraph).
		Handle("pre", handleChildren).
		Handle("table", c.table).
		Handle("small", handleChildren).
		Handle("span", handleChildren).
		Handle("strong", handleChildren).
		Handle("sup", handleChildren).
		Handle("ul", c.ul).
		// These elements are ignored
		Handle("bookMeta", nil)

	// Final Handler
	c.handler = parser.Of(
		c.beginDocument,
		parser.New().
			Handle("div",
				content.Handler().
					HasClasses("td-content", "td-content-new").
					HandleChildren).
			Default(handleChildren).
			Handler().
			FindByClass("td-main", "td-main-new"),
		c.endDocument)

	return c, nil
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
