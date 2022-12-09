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
	Styles map[string][]*Style `yaml:"styles"`
}

type Style struct {
	Rule    *Rule             `yaml:"-"`    // The compiled rule
	RuleSrc string            `yaml:"rule"` // Rule e.g. td:not(:last-child)
	Css     map[string]string `yaml:"css"`
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
func (c *Css) Start() error {
	log.Println("Loading LaTeX CSS stylesheets")
	for k, v := range c.Styles.Styles {
		log.Println(k)
		for _, ve := range v {
			log.Println(ve.RuleSrc, ve.Css)
			r, err := ParseRule(ve.RuleSrc)
			if err != nil {
				return err
			}
			ve.Rule = r
		}
	}

	for _, v := range c.Styles.Styles {
		for _, ve := range v {
			ve.Rule.Write(os.Stdout)
		}
	}

	return io.EOF
}
