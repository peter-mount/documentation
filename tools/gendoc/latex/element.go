package latex

import (
	"github.com/peter-mount/documentation/tools/gendoc/latex/parser"
	"github.com/peter-mount/documentation/tools/gendoc/latex/table"
	"github.com/peter-mount/documentation/tools/gendoc/latex/util"
	"golang.org/x/net/html"
)

func (l *LaTeX) parseChildren(w *util.Writer, n *html.Node) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := l.parseNode(w, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LaTeX) parseNode(w *util.Writer, n *html.Node) error {
	switch n.Type {
	case html.ElementNode:
		err := l.parseElement(w, n)
		if err != nil {
			if parser.IsStopChildTraverse(err) {
				return nil
			}
			return err
		}
		/*
		   case html.TextNode:
		     s := strings.TrimSpace(n.Data)
		     if len(s) > 0 {
		       _ = w.WriteString("%s\n", strings.ReplaceAll(s, "&", "\\&"))
		     }*/

	}

	return l.parseChildren(w, n)
}

func (l *LaTeX) parseElement(w *util.Writer, n *html.Node) error {
	// Don't parse elements with this class as they are not for print
	if util.CheckClass(n, "d-print-none") {
		return parser.StopChildTraverse()
	}

	switch n.DataAtom.String() {
	case "h1":
		//w.WriteString(util.Command2, "title ")
		err := l.parseChildren(w, n)
		if err != nil {
			return err
		}
		//w.WriteString(util.Command3)
		return parser.StopChildTraverse()

	case "table":
		t := table.Table{}

		err := t.Parse(n)
		if err != nil {
			return err
		}

		err = t.Write(w)
		if err != nil {
			return err
		}

		return parser.StopChildTraverse()
	}

	return nil
}
