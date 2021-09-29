package html

import (
  "github.com/peter-mount/documentation/tools/util"
  "go/token"
)

// ChipExpansion takes a chip pin definition and converts it into HTML
func (e *Element) ChipExpansion(s string) *Element {
  return e.chipExpansion(util.Tokenize(s))
}

func (e *Element) chipExpansion(t *util.Token) *Element {
  ce := e
  for t != nil {
    if t.Child != nil {
      switch t.Token {
      case token.SEMICOLON:
      // Ignore
      case token.LPAREN:
        ce = ce.Text(" (").chipExpansion(t.Child).Text(")")
      case token.IDENT:
        switch t.Lit {
        case "A":
          ce = ce.Text("A").Sub().chipExpansion(t.Child).End()
        case "D":
          ce = ce.Text("D").Sub().chipExpansion(t.Child).End()
        case "P":
          ce = ce.Text("P").Sub().chipExpansion(t.Child).End()
        case "V":
          ce = ce.Text("V").Sub().chipExpansion(t.Child).End()
        case "NOT":
          ce = ce.Span().Class("not").chipExpansion(t.Child).End()
        case "PHI":
          ce = ce.Text("∅").Sub().chipExpansion(t.Child).End()
        case "SUB":
          ce = ce.Sub().chipExpansion(t.Child).End()
        case "SUP":
          ce = ce.Sup().chipExpansion(t.Child).End()
        default:
          ce = ce.Text(t.String()).chipExpansion(t.Child)
        }
      default:
        ce = ce.Text(t.String()).chipExpansion(t.Child)
      }
    } else if t.Token != token.SEMICOLON {
      ce = ce.Text(t.String())
    }
    t = t.Next
  }
  return ce
}
