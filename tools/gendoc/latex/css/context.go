package css

import "golang.org/x/net/html"

type Context struct {
	Doc     *html.Node
	Node    *html.Node
	matches []*html.Node // Slice of matched nodes
}

func (c *Context) Reset() {
	c.matches = nil
}

// AddMatch adds a node to the matches list.
// Note: This function will ensure that a node only exists in the slice once.
func (c *Context) AddMatch(n *html.Node) *Context {
	for _, e := range c.matches {
		if e == n {
			return c
		}
	}
	c.matches = append(c.matches, n)
	return c
}

func (c *Context) HasMatches() bool {
	return len(c.matches) > 0
}

func (c *Context) ForEachMatch(f func(*html.Node) error) error {
	for _, n := range c.matches {
		err := f(n)
		if err != nil {
			return err
		}
	}
	return nil
}
