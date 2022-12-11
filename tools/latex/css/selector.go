package css

import (
	"github.com/peter-mount/documentation/tools/latex/parser"
	"golang.org/x/net/html"
)

func (s *Styles) Search(doc *html.Node) error {
	return search(&Context{Doc: doc}, doc, s.Styles)
}

func search(ctx *Context, doc *html.Node, styles []*Style) error {
	return parser.TraverseChildren(doc, 0, func(node *html.Node) error {
		for _, style := range styles {
			ctx.Node = node
			ctx.Reset()

			err := style.search(ctx)
			if err == nil && ctx.HasMatches() {
				err = ctx.ForEachMatch(style.Apply)

				// If we have a match then process any children
				if err == nil && len(style.Children) > 0 {
					err = search(ctx, node, style.Children)
				}
			}

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Style) search(ctx *Context) error {
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
