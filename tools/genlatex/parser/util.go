package parser

import "golang.org/x/net/html"

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
