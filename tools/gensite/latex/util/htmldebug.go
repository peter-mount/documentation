package util

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
)

// DebugHtml will write a simple tree diagram of a Node to a Writer
func DebugHtml(n *html.Node, w io.Writer) error {
	return debugHtml("", n, w)
}

func debugHtml(ps string, n *html.Node, w io.Writer) error {
	var s, ns string
	if n.NextSibling == nil {
		s = ps + "+-"
		ns = ps + "  "
	} else {
		s = ps + "|-"
		ns = ps + "| "
	}

	var err error
	switch n.Type {
	case html.DocumentNode:
		_, err = fmt.Fprintf(w, "%sDocument\n", s)
	case html.ElementNode:
		_, err = fmt.Fprintf(w, "%sElem %q\n", s, n.Data)
	case html.TextNode:
		_, err = fmt.Fprintf(w, "%sText %q\n", s, n.Data)
	default:
		_, err = fmt.Fprintf(w, "%sUnk%v %q\n", s, n.Type, n.Data)
	}
	if err != nil {
		return err
	}

	for _, attr := range n.Attr {
		_, err := fmt.Fprintf(w, "%sAttr %q=%q\n", s, attr.Key, attr.Val)
		if err != nil {
			return err
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = debugHtml(ns, c, w)
		if err != nil {
			return err
		}
	}

	return nil
}
