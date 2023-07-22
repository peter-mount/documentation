package latex

import (
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"gopkg.in/yaml.v2"
	"os"
)

type Converter struct {
	handler    parser.Handler
	stylesheet *stylesheet.Stylesheet
}

func (c *Converter) Handler() parser.Handler {
	return c.handler
}

func (c *Converter) Stylesheet() *stylesheet.Stylesheet {
	return c.stylesheet
}

func New(config string) (*Converter, error) {
	c := &Converter{}

	// load stylesheet
	b, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	c.stylesheet = &stylesheet.Stylesheet{}
	err = yaml.Unmarshal(b, &c.stylesheet)
	if err != nil {
		return nil, err
	}

	// Ensure it's fully configured
	c.stylesheet.Init()

	// Handle normal html content
	content := parser.New()

	// This handler is used by multiple elements
	heading := parser.Of(c.headingStart, content.Handler().HandleChildren, c.headingEnd)

	content.
		Text(c.text).
		// anchor is ignored, just parse it's content
		Handle("a", util.HandleChildren).
		Handle("br", c.lineBreak).
		Handle("code", c.code).
		Handle("div", c.div).
		Handle("dd", util.HandleChildren).
		Handle("dl", util.HandleChildren).
		Handle("dt", util.HandleChildren).
		Handle("em", c.em).
		Handle("figure", parser.Of(
			c.figureStart,
			parser.New().
				Handle("figcaption", nil).
				//Handle("figcaption", parser.Of(c.figureCaptionStart, content.Handler().HandleChildren, c.figureCaptionEnd)).
				Default(content.Handler()).
				Handler().
				HandleChildren,
			c.figureEnd1,
			//contentHandler.Type("figcaption").HandleChildren,
			c.figureEnd2)).
		Handle("h1", heading).
		Handle("h2", heading).
		Handle("h3", heading).
		Handle("h4", heading).
		Handle("h5", heading).
		Handle("li", c.li).
		Handle("ol", c.ol).
		Handle("p", c.paragraph).
		Handle("pre", util.HandleChildren).
		Handle("table", c.table).
		Handle("small", util.HandleChildren).
		Handle("span", util.HandleChildren).
		Handle("strong", c.strong).
		Handle("sup", c.sup).
		Handle("ul", c.ul).
		// These elements are ignored
		Handle("bookMeta", nil)

	// Final Handler
	c.handler = parser.Of(
		c.beginDocument,
		parser.New().
			Handle("div",
				content.Handler().
					HasClasses("td-content", "td-content-new").
					HandleChildren).
			Default(util.HandleChildren).
			Handler().
			FindByClass("td-main", "td-main-new"),
		c.endDocument)

	return c, nil
}
