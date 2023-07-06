package parser

import (
	"golang.org/x/net/html"
	"strings"
)

func FindNode(n *html.Node, p Predicate) *html.Node {
	if p(n) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := FindNode(c, p); r != nil {
			return r
		}
	}
	return nil
}

func FindByClass(n *html.Node, classes ...string) *html.Node {
	return FindNode(n, NewPredicate().HasClass(classes...))
}

func FindById(n *html.Node, id string) *html.Node {
	return FindNode(n, func(n *html.Node) bool {
		return GetAttr(n, "id") == id
	})
}

func GetText(n *html.Node) string {
	var s []string
	s = getText(s, n)
	return strings.Join(s, " ")
}

func getText(s []string, n *html.Node) []string {
	switch n.Type {
	case html.TextNode:
		s = append(s, n.Data)
	case html.ElementNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s = getText(s, c)
		}
	}
	return s
}
