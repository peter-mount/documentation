package css

import (
	"github.com/peter-mount/documentation/tools/latex/parser"
	"golang.org/x/net/html"
)

const (
	keyDoc  = "doc"
	keyNode = "node"
)

func (s *Styles) Search(doc *html.Node) error {
	ctx := &Context{Doc: doc}

	return parser.TraverseChildren(doc, 0, func(node *html.Node) error {
		ctx.Node = node

		for _, styles := range s.Styles {
			for _, style := range styles {
				ctx.Reset()

				err := style.Search(ctx)
				if err != nil {
					return err
				}

				if ctx.HasMatches() {
					err = ctx.ForEachMatch(style.Apply)
				}
			}
		}

		return nil
	})
}

func (s *Style) Search(ctx *Context) error {
	for _, nodes := range s.Rule.Nodes {
		v, err := nodes.Do(ctx)
		if err != nil {
			return err
		}

		if v.Bool() {
			ctx.AddMatch(ctx.Node)
		}
	}
	return nil
}
