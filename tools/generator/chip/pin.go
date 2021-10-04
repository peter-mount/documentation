package chip

import "github.com/peter-mount/documentation/tools/util/html"

func dipPinLabel2(align bool, pin, x, y, r int, l string, e *html.Element) *html.Element {
  z := 8 + dipPinWidth
  px := -1
  tx := dipPinWidth2
  if align {
    z = dipPinWidth - z
    tx = tx + 4
  } else {
    px = dipPinWidth - 8 + px
    tx = tx - 6
  }

  e = e.G().Attr("transform", "translate(%d %d) rotate(%d)", x, y, -90*r).
    Rect().Class("chipPin").X(px).End().
    SvgText().Class("chipPin").X(tx).Y(dipPinHeight2+6).Textf("%d", pin).End()
  e = dipPinLabel(align, z, 17, l, e)
  return e.End()
}

func dipPinLabel(rightAlign bool, x, y int, label string, e *html.Element) *html.Element {
  b := e.G().Attr("transform", "translate(%d %d)", x, y)
  _ = Parse(label, dipPinFontSize, rightAlign, b)
  return b.End()
}
