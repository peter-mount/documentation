package chip

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util/html"
  "github.com/peter-mount/go-kernel/util"
  "strings"
)

const (
  // Valid row labels, no I O or Q as could be read as 1, 0 or 0 respectively
  // TODO verify if O & Q are excluded, this is a presumption but have seen I
  pgaRowLabels = "ABCDEFGHJKLMNPRSTUVWXYZ"
)

// pga Pin Grid Array
//
// pinCount is the number of pins on both axes. Not all pins need to be defined.
// Pins are labeled A1, A2, K10 etc where A..Z is Y-axis whilst 1..n X-axis
func pga(d *Definition) error {
  b := html.Builder()

  spacing := dipPinSpacingV

  height := spacing * (d.PinCount + 1)
  height2 := height / 2
  width := height // These are square chips
  width2 := width / 2

  vWidth := 50 + width // TODO default at 200, make configurable to increase the size
  vWidth2 := vWidth / 2
  vHeight := vWidth
  vHeight2 := vHeight / 2

  // center chip die size
  bSize := 50

  b = b.
    Svg().ViewBox(0, 0, vWidth, vHeight).Width(vWidth).
    Style().Attr("type", "text/css").
    // TODO check text-decoration-thickness as seems broken in Chrome 94 but works in Firefox 92
    Text("tspan.not {text-decoration: overline;webkit-text-decoration-thickness: 2px;text-decoration-thickness: 2px;}").
    Textf(".chip {font-family:%s;}", html.TextFont).
    Textf("g.chipCase polygon {stroke:black;fill:none;stroke-width: 2px;}").
    Textf("g.chipCase line {stroke:black;fill:none;stroke-width: 2px;}").
    Text(".chipCase circle {fill:white;}").
    Text(".chipDie {opacity:30%;}").
    Text(".chipDie .chipDieDash {stroke-dasharray:4 4;}").
    Text("g.chipLabel {fill:black;text-anchor:middle;opacity:30%;}").
    Textf("text.chipLabel {font-size:%dpx;font-weight:bold;}", 20).
    Textf("text.chipSubLabel {font-size:%dpx;}", 16).
    Textf("text.chipTypeLabel {font-size:%dpx;}", 10).
    Text(".pins circle {fill:none;stroke:black;strike-width: 1px;}").
    Textf(".pins text {fill:black;text-anchor:middle;font-size:%dpx;}", 10).
    End(). // style
    G().Class("chip").
    // DIP visualisation - last so it overlays the pins
    G().Class("chipCase").
    Polygon().
    Point(vWidth2-width2+16, vHeight2-height2).
    Point(vWidth2-width2, vHeight2-height2+16).
    Point(vWidth2-width2, vHeight2+height2-16).
    Point(vWidth2-width2+16, vHeight2+height2).
    Point(vWidth2+width2-16, vHeight2+height2).
    Point(vWidth2+width2, vHeight2+height2-16).
    Point(vWidth2+width2, vHeight2-height2+16).
    Point(vWidth2+width2-16, vHeight2-height2).
    End(). //Polygon()
    G().Class("chipDie").
    Polygon().
    Point(vWidth2-bSize, vHeight2-bSize).
    Point(vWidth2-bSize, vHeight2+bSize).
    Point(vWidth2+bSize, vHeight2+bSize).
    Point(vWidth2+bSize, vHeight2-bSize).
    End(). //Polygon()
    Line().X1(vWidth2-width2).Y1(vHeight2+height2-bSize).X2(vWidth2-width2+bSize).Y2(vHeight2+height2-bSize).End().
    Line().X1(vWidth2-width2+bSize).Y1(vHeight2+height2-bSize).X2(vWidth2-width2+bSize).Y2(vHeight2+height2).End().
    Line().Class("chipDieDash").X1(vWidth2-width2+bSize).Y1(vHeight2+height2-bSize).X2(vWidth2-bSize).Y2(vHeight2+bSize).End().
    End(). // G.chipDie
    End(). // G.chipCase
    // DIP Labels
    G().Class("chipLabel").
    Attr("transform", "translate(%d %d)", vWidth2, vHeight2).
    SvgText().Y(-27).Class("chipLabel").Text(d.Label).End().
    SvgText().Y(0).Class("chipSubLabel").Text(d.SubLabel).End().
    SvgText().Y(30).Class("chipTypeLabel").Text(strings.ToUpper(d.Type)).End().
    SvgText().Y(40).Class("chipTypeLabel").Text("Bottom View").End().
    End(). // G.chipLabel
    // Pin Labels
    G().Class("pinLabel").
      Sequence(1, d.PinCount, func(pin int, e *html.Element) *html.Element {
        ofs := 15
        xy := spacing * pin

        e = e.G().Attr("transform", "translate(%d %d)", 0, xy).
          SvgText().Class("chipPin").X(ofs - 10).Y(dipPinHeight2 + 15).
          Text(string(pgaRowLabels[d.PinCount-pin])).
          End().
          End()

        e = e.G().Attr("transform", "translate(%d %d)", xy, 0).
          SvgText().Class("chipPin").X(dipPinWidth2+6).Y(vHeight-ofs+10).
          Textf("%d", pin).
          End().
          End()

        return e
      }).
    End(). // G.pinLabel
    G().Class("pins").
      Sequence(1, d.PinCount, func(y int, e *html.Element) *html.Element {
        return e.Sequence(1, d.PinCount, func(x int, e *html.Element) *html.Element {
          if l, exists := d.Pins[x-1+((y-1)*d.PinCount)]; exists {
            px := (spacing * x) + 28
            py := (spacing * (d.PinCount - y + 1)) + 13 + 9

            e = e.Circle().CX(px).CY(py).R(5).End()

            te := e.G().Attr("transform", "translate(%d %d)", px, py+18)
            _ = Parse(l, 9, false, te)
            e = te.End()
          }

          return e
        })
      }).
    End() // pins

  // =================================================================
  // Convert to FileBuilder & generate the shortcode html file
  return b.
    End(). // Svg
    FileBuilder().
    FileHandler().
    Write(d.Path("static/static/chipref")+".svg", d.FileInfo.ModTime())
}

// pgaDecodePins decodes Pins map in yaml to the matrix for display
func pgaDecodePins(v *Definition, m map[interface{}]interface{}) error {
  for x := 0; x < v.PinCount; x++ {
    for y := 0; y < v.PinCount; y++ {
      k := fmt.Sprintf("%c%d", pgaRowLabels[y], x+1)
      if l, exists := m[k]; exists {
        v.Pins[x+(v.PinCount*y)] = util.DecodeString(l, "??")
      }
    }
  }
  return nil
}
