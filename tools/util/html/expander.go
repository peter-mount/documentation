package html

import "strings"

func (e *Element) ChipExpansion(s string) *Element {
  for s != "" {
    is := strings.Index(s, "(")
    ie := strings.Index(s, ")")
    if ie > is && is > -1 {
      switch s[:is] {
      // Address line
      case "A":
        e.Text("A").Sub().ChipExpansion(s[is+1 : ie]).End().End()
      // Data line
      case "D":
        e.Text("D").Sub().ChipExpansion(s[is+1 : ie]).End().End()
      // Peripheral I/O line
      case "P":
        e.Text("P").Sub().ChipExpansion(s[is+1 : ie]).End().End()
      case "NOT":
        e.Span().Class("not").ChipExpansion(s[is+1 : ie]).End()
      case "PHI":
        e.Text("∅").Sub().ChipExpansion(s[is+1 : ie]).End().End()
      case "SUB":
        e.Sub().ChipExpansion(s[is+1 : ie]).End().End()
      case "SUP":
        e.Sup().ChipExpansion(s[is+1 : ie]).End().End()
      default:
        e.Text(s[:ie+1])
      }
      s = s[ie+1:]
    } else if ie > -1 {
      e.Text(s[:ie+1])
      s = s[ie+1:]
    } else {
      e.Text(s)
      s = ""
    }
  }
  return e
}
