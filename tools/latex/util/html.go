package util

import (
	"golang.org/x/net/html"
	"strings"
)

type HtmlHandler func(*Writer, *html.Node) error

// WriteHtml handles converting standard html elements to Html
func (w *Writer) WriteHtml(n *html.Node) error {

	switch n.Type {
	case html.TextNode:
		s := strings.TrimSpace(n.Data)
		if len(s) > 0 {
			_ = w.WriteString(EscapeText(s))
		}

	case html.ElementNode:
		h, exists := elementHandlers[n.DataAtom.String()]
		if exists {
			return h(w, n)
		}
	}

	return nil
}

var elementHandlers = map[string]HtmlHandler{
	"br": newLine,
	// Table internals are ignored as handled by Table
	"thead": nop, "tbody": nop, "tr": nop, "th": nop, "td": nop, "caption": nop,
}

func nop(w *Writer, _ *html.Node) error {
	return nil
}

func newLine(w *Writer, _ *html.Node) error {
	_ = w.WriteString("\n")
	return nil
}
