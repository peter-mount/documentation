package parser

import (
	"golang.org/x/net/html"
)

type Predicate func(*html.Node) bool

func (a Predicate) Do(n *html.Node) bool {
	if a == nil {
		return false
	}
	return a(n)
}

func True(_ *html.Node) bool { return true }

func False(_ *html.Node) bool { return false }

func NewPredicate() Predicate { return nil }

func (a Predicate) And(b Predicate) Predicate {
	if a == nil || b == nil {
		return False
	}
	return func(n *html.Node) bool {
		return a(n) && b(n)
	}
}

func (a Predicate) Or(b Predicate) Predicate {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(n *html.Node) bool {
		return a(n) || b(n)
	}
}

func (a Predicate) Not() Predicate {
	if a == nil {
		return False
	}
	return func(n *html.Node) bool {
		return !a(n)
	}
}

func (a Predicate) HasClass(classes ...string) Predicate {
	if len(classes) == 0 {
		return False
	}

	return a.Or(func(n *html.Node) bool {
		nodeClasses := GetClass(n)
		if len(nodeClasses) > 0 {
			for _, c := range classes {
				for _, e := range nodeClasses {
					if e == c {
						return true
					}
				}
			}
		}
		return false
	})
}
