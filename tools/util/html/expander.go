package html

import (
  "github.com/peter-mount/documentation/tools/util"
  "go/token"
)

// ChipExpansion takes a chip pin definition and converts it into HTML.
//
// This is a simple expression format.
// To put a line above text then use NOT(text).
// For subscript use SUB(text).
// For superscript use SUP(text).
//
// e.g. "A SUB(10)" will generate html for A10 where 10 is subscript.
//
// Common labels have shortcuts, so instead of "A SUB(10)" you can use "A(10)" to do the same thing.
// Shortcuts provided are: A, D, P, V & PHI (reality the unicode open set symbol but looks like the greek Phi character)
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
        // Text with a line above it
        case "NOT":
          ce = ce.Span().Class("not").chipExpansion(t.Child).End()
        // Subscript text
        case "SUB":
          ce = ce.Sub().chipExpansion(t.Child).End()
        // Superscript text
        case "SUP":
          ce = ce.Sup().chipExpansion(t.Child).End()
        // Common shortcuts for "id SUB(val)" so we can use "id(val)" instead
        case "A":
          ce = ce.Text("A").Sub().chipExpansion(t.Child).End()
        case "D":
          ce = ce.Text("D").Sub().chipExpansion(t.Child).End()
        case "P":
          ce = ce.Text("P").Sub().chipExpansion(t.Child).End()
        case "V":
          ce = ce.Text("V").Sub().chipExpansion(t.Child).End()
        case "PHI":
          ce = ce.Text("∅").Sub().chipExpansion(t.Child).End()
        // Default the ident as text then expand the child content
        default:
          ce = ce.Text(t.String()).chipExpansion(t.Child)
        }
      // Default the ident as text then expand the child content
      default:
        ce = ce.Text(t.String()).chipExpansion(t.Child)
      }
    } else if t.Token != token.SEMICOLON {
      // For all content other than token.SEMICOLON just append it as text
      ce = ce.Text(t.String())
    }
    t = t.Next
  }
  return ce
}
