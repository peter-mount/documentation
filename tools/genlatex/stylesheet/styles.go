package stylesheet

import (
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

type Style struct {
	FontSize      string `yaml:"fontSize"`      // LaTeX font size, e.g. \normalsize \footnotesize
	Align         string `yaml:"align"`         // Alignment of cell default c
	VerticalAlign string `yaml:"verticalAlign"` // Vertical Align, can be t, c or b
}

// GetStyle returns a Style based on the classes on a node.
// This will return a composite, so the first entries will override latter ones
func (s *Stylesheet) GetStyle(n *html.Node) Style {
	r := Style{}
	for _, c := range parser.GetClass(n) {
		if t, exists := s.Styles[c]; exists {
			r.Align = DefString(r.Align, t.Align)
			r.VerticalAlign = DefString(r.VerticalAlign, t.VerticalAlign)
		}
	}
	return r
}
