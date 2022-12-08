package parser

import (
	"golang.org/x/net/html"
	"strings"
)

// GetAttribute returns the named attribute on a node, if it exists
func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func GetId(n *html.Node) (string, bool) {
	return GetAttribute(n, "id")
}

// CheckId returns true if the node is an element, and it's id attribute matches the supplied id.
func CheckId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetId(n)
		return ok && s == id
	}
	return false
}

func GetClass(n *html.Node) (string, bool) {
	return GetAttribute(n, "class")
}

// CheckClass returns true if a node is an element, and it has a specific class
func CheckClass(n *html.Node, class string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetClass(n)
		if ok {
			if s == class {
				return true
			}
			if strings.Contains(s, class) {
				// We need to ensure the class is a specific word not just it exists
				for _, c := range strings.Split(s, " ") {
					if c == class {
						return true
					}
				}
			}
		}
	}
	return false
}
