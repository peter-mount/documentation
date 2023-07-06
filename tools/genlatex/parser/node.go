package parser

import (
	"context"
	"golang.org/x/net/html"
)

type Handler func(n *html.Node, ctx context.Context) error

func Of(handlers ...Handler) Handler {
	var r Handler
	for _, h := range handlers {
		r = r.Then(h)
	}
	return r
}

// Then returns a Handler that will call this Handler first then, if no error has occurred
// the second Handler.
func (a Handler) Then(b Handler) Handler {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(n *html.Node, ctx context.Context) error {
		err := a(n, ctx)
		if err == nil {
			err = b(n, ctx)
		}
		return err
	}
}

// Do will call a handler.
// If the Handler is nil this does nothing.
func (a Handler) Do(n *html.Node, ctx context.Context) error {
	if a != nil {
		return a(n, ctx)
	}
	return nil
}

// HandleChildren will call the Handler for each child in a html.Node
//
// If the Handler is nil this does nothing.
func (a Handler) HandleChildren(n *html.Node, ctx context.Context) error {
	if a == nil {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := a(c, ctx); err != nil {
			return err
		}
	}
	return nil
}

// HandleChildren will call a Handler for each child in a html.Node
//
// If the Handler is nil this does nothing.
func HandleChildren(h Handler, n *html.Node, ctx context.Context) error {
	return h.HandleChildren(n, ctx)
}

// Type returns a Handler that will call this handler only
// if the node being passed is a html.ElementNode and it's of this type.
// e.g. "div"
//
// If the Handler is nil this does nothing.
func (a Handler) Type(t string) Handler {
	if a == nil {
		return nil
	}
	return func(n *html.Node, ctx context.Context) error {
		if n.Type == html.ElementNode && n.Data == t {
			return a(n, ctx)
		}
		return nil
	}
}

// HasClass returns a Handler that will call this handler only
// if the node being passed is a html.ElementNode, it has the class
// attribute with the specified class defined.
//
// If the Handler is nil this does nothing.
func (a Handler) HasClass(c string) Handler {
	if a == nil {
		return nil
	}
	return func(n *html.Node, ctx context.Context) error {
		if HasClass(n, c) {
			return a(n, ctx)
		}
		return nil
	}
}

// HasClasses will invoke this Handler if a node has one of these classes.
//
// If Handler is nil or classes is empty then this returns nil.
//
// If classes is a single entry this is the same as HasClass.
func (a Handler) HasClasses(classes ...string) Handler {
	if a == nil {
		return nil
	}

	switch len(classes) {
	case 0:
		// Nothing to test so return nil as we will never invoke
		return nil

	case 1:
		// HasClass is more performant for a single class
		return a.HasClass(classes[0])

	default:
		return a.If(NewPredicate().HasClass(classes...))
	}
}

// If will invoke this Handler if a Predicate passes.
//
// If either Handler or Predicate is nil this will return nil
func (a Handler) If(p Predicate) Handler {
	if a == nil || p == nil {
		return nil
	}
	return func(n *html.Node, ctx context.Context) error {
		if p(n) {
			return a(n, ctx)
		}
		return nil
	}
}

func (a Handler) FindByClass(classes ...string) Handler {
	if a == nil {
		return nil
	}

	return func(n *html.Node, ctx context.Context) error {
		n = FindByClass(n, classes...)
		if n != nil {
			return a(n, ctx)
		}
		return nil
	}
}
