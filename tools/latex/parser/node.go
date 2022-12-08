package parser

import (
	"errors"
	"golang.org/x/net/html"
)

var (
	stopTraverse      = errors.New("traverse stopped")
	stopChildTraverse = errors.New("traverse children stopped")
)

func IsStopTraverse(err error) bool {
	return err == stopTraverse
}

func StopTraverse() error {
	return stopTraverse
}

func IsStopChildTraverse(err error) bool {
	return err == stopChildTraverse
}

func StopChildTraverse() error {
	return stopChildTraverse
}

// Traverse traverses a node tree calling a function for each node found.
// If t is not 0 then it's a NodeType to match before calling the function.
func Traverse(n *html.Node, t html.NodeType, f func(*html.Node) error) error {
	if t == 0 || t == n.Type {
		err := f(n)
		if err != nil {
			if IsStopChildTraverse(err) {
				// Return to caller but carry on, just don't run children
				return nil
			}
			return err
		}
	}

	return TraverseChildren(n, t, f)
}

func TraverseChildren(n *html.Node, t html.NodeType, f func(*html.Node) error) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := Traverse(c, t, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// ForEachChild will call a function just for the children of a node.
// If t is not 0 then it's a NodeType to match before calling the function.
func ForEachChild(n *html.Node, t html.NodeType, f func(*html.Node) error) error {
	return ForEachSibling(n.FirstChild, t, f)
}

// ForEachSibling will call a function just for the sibling of a node, starting with
// the node.
// If t is not 0 then it's a NodeType to match before calling the function.
func ForEachSibling(n *html.Node, t html.NodeType, f func(*html.Node) error) error {
	for c := n; c != nil; c = c.NextSibling {
		if t == 0 || t == n.Type {
			err := f(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
