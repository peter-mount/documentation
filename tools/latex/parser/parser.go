package parser

import (
	"golang.org/x/net/html"
	"net/http"
)

// Parser represents a html document being processed into LaTeX
type Parser struct {
	url   string                // Base url for the html, used for retrieving additional resources
	doc   *html.Node            // Document root node
	head  *html.Node            // Html Head node
	body  *html.Node            // Html Body node
	idMap map[string]*html.Node // Map of id's to node, to speed up lookups
}

// Parse will retrieve a document from a URL and return a parser for it
func Parse(url string) (*Parser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	parser := &Parser{
		url:   url,
		doc:   doc,
		idMap: make(map[string]*html.Node),
	}

	// For efficiency, locate the head and body elements
	_ = ForEachSibling(doc, html.DocumentNode, func(n *html.Node) error {
		return ForEachChild(n, html.DoctypeNode, func(n *html.Node) error {
			return ForEachChild(n, html.ElementNode, func(n *html.Node) error {
				switch n.Data {
				case "head":
					parser.head = n
				case "body":
					parser.body = n
				}
				return nil
			})
		})
	})

	// Init idMap starting from the body element
	_ = Traverse(parser.body, html.ElementNode, func(n *html.Node) error {
		if id, exists := GetId(n); exists {
			if _, exists = parser.idMap[id]; !exists {
				parser.idMap[id] = n
			}
		}
		return nil
	})

	return parser, nil
}

// GetElementById returns a node by its id
func (p *Parser) GetElementById(id string) *html.Node {
	return p.idMap[id]
}

// GetElementByClass returns a slice of nodes based on the class name
// There's no caching here like GetElementById, so it's expensive each time it's called
func (p *Parser) GetElementByClass(class string) []*html.Node {
	var ret []*html.Node

	_ = Traverse(p.body, html.ElementNode, func(n *html.Node) error {
		if CheckClass(n, class) {
			ret = append(ret, n)
		}
		return nil
	})

	return ret
}
