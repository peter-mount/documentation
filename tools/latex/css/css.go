package css

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
)

// Css represents a simple CSS stylesheet for rendering latex
type Css struct {
	Styles *Styles `kernel:"config,css"`
}

type Styles struct {
	Styles []*Style `yaml:"styles"`
}

type Style struct {
	Rule     *Rule             `yaml:"-"`        // The compiled rule
	RuleSrc  string            `yaml:"rule"`     // Rule e.g. td:not(:last-child)
	Css      map[string]string `yaml:"css"`      // CSS to apply
	Children []*Style          `yaml:"children"` // Child rules, processed only if this one matches
}

func (s *Style) Apply(n *html.Node) error {
	for k, v := range s.Css {
		s.ApplyCSS(n, k, v)
	}
	return nil
}

func (s *Style) ApplyCSS(n *html.Node, k, v string) {
	key := "css-" + k

	// Overwrite any existing entry
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = v
			return
		}
	}

	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: v})
}

func Parse(s []*Style) error {
	for _, ve := range s {
		r, err := ParseRule(ve.RuleSrc)
		if err == nil {
			ve.Rule = r
			if len(ve.Children) > 0 {
				err = Parse(ve.Children)
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func Write(w io.Writer, s []*Style) {
	for _, ve := range s {
		ve.Write(w)
	}
}

func (s *Style) Write(w io.Writer) {
	s.Rule.Write(w)
}

func (c *Css) Start() error {
	log.Println("Loading LaTeX CSS stylesheets")
	err := Parse(c.Styles.Styles)
	if err != nil {
		return err
	}

	Write(os.Stdout, c.Styles.Styles)

	return io.EOF
}
