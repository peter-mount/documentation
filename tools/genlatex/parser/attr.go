package parser

import (
	"golang.org/x/net/html"
	"strings"
)

// GetAttr returns the named attribute in a html.Node
func GetAttr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// HasClass returns true if a html.Node has a class attribute and it has
// the specified class name declared within it.
func HasClass(n *html.Node, name string) bool {
	if n.Type == html.ElementNode {
		a := GetAttr(n, "class")
		if a != "" {
			for _, e := range strings.Split(a, " ") {
				if e == name {
					return true
				}
			}
		}
	}
	return false
}
