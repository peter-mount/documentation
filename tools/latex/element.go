package latex

import (
  "github.com/peter-mount/documentation/tools/latex/parser"
  "golang.org/x/net/html"
  "log"
  "strings"
)

func (l *LaTeX) parseChildren(w *Writer, n *html.Node) error {
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    err := l.parseNode(w, c)
    if err != nil {
      return err
    }
  }
  return nil
}

func (l *LaTeX) parseNode(w *Writer, n *html.Node) error {
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

func (l *LaTeX) parseElement(w *Writer, n *html.Node) error {
  // Don't parse elements with this class as they are not for print
  if parser.CheckClass(n, "d-print-none") {
    return parser.StopChildTraverse()
  }

  switch n.DataAtom.String() {
  case "h1":
    w.WriteString(Command2, "title ")
    err := l.parseChildren(w, n)
    if err != nil {
      return err
    }
    w.WriteString(Command3)
    return parser.StopChildTraverse()

  case "table":
    t := NewTable(w, n)
    log.Printf("table %d %d", t.maxCols, t.rows)

    w = w.Begin("center").Begin("tabular").
      WriteString("{|| %s||}", strings.Repeat("c ", t.maxCols))

    err := parser.TraverseChildren(n, 0, t.write)
    if err != nil {
      return err
    }

    w.WriteString("\\hline\n")
    w = w.End().End()
    return parser.StopChildTraverse()
  }

  return nil
}
