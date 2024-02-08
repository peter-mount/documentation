package parser

import (
	"golang.org/x/net/html"
	"strconv"
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

func GetAttrInt(n *html.Node, name string, def int) (int, bool) {
	a := GetAttr(n, name)
	if a != "" {
		i, err := strconv.Atoi(a)
		if err == nil {
			return i, true
		}
	}

	return def, false
}

// HasClass returns true if a html.Node has a class attribute and it has
// the specified class name declared within it.
func HasClass(n *html.Node, name string) bool {
	classes := GetClass(n)
	if len(classes) > 0 {
		for _, e := range classes {
			if e == name {
				return true
			}
		}
	}
	return false
}

func HasClasses(n *html.Node, required ...string) bool {
	classes := GetClass(n)
	if len(classes) > 0 {
		for _, r := range required {
			for _, e := range classes {
				if e == r {
					return true
				}
			}
		}
	}
	return false
}

func GetClass(n *html.Node) []string {
	if n.Type == html.ElementNode {
		a := GetAttr(n, "class")
		if a != "" {
			return strings.Split(a, " ")
		}
	}
	return nil
}
