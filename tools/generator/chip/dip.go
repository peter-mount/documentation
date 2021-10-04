package chip

import (
  "github.com/peter-mount/documentation/tools/util/html"
)

const (
  dipLabelFontSize    = 50
  dipSubLabelFontSize = 32
  dipPinFontSize      = 14
  dipPinWidth         = 32
  dipPinWidth2        = dipPinWidth / 2
  dipPinHeight        = 22
  dipPinHeight2       = dipPinHeight / 2
  dipPinSpacingV      = dipPinHeight + 8
  dipPinSpacingV2     = dipPinSpacingV / 2
  dipWidth            = 125
  dipWidth2           = dipWidth / 2
)

// dip Dual Inline Pin chip layout
func dip(d *Definition) error {
  b := html.HtmlBuilder()

  pinCount2 := d.PinCount / 2

  height := (dipPinSpacingV * pinCount2) + 20
  height2 := height / 2
  vWidth := 300 + dipWidth
  vWidth2 := vWidth / 2
  b.
    Svg().ViewBox(0, 0, vWidth, height+50).Width(vWidth).
    Style().Attr("type", "text/css").
    // TODO check text-decoration-thickness as seems broken in Chrome 94 but works in Firefox 92
    Text("tspan.not {text-decoration: overline;webkit-text-decoration-thickness: 2px;text-decoration-thickness: 2px;}").
    Textf(".chip {font-family:%s;}", html.TextFont).
    Textf("g.chipCase rect {width:%dpx;height:%dpx;fill:#333;}", dipWidth, height).
    Text(".chipCase circle {fill:white;}").
    Textf("rect.chipPin {width:%dpx;height:%dpx;fill:%s;stroke:%s;}", dipPinWidth, dipPinHeight, html.LightGrey, html.BLACK).
    Textf("text.chipPin {font-size:%dpx;text-anchor:middle;}", dipPinFontSize).
    Textf("g.chipLabel {fill:lightgrey;text-anchor:middle;}").
    Textf("text.chipLabel {font-size:%dpx;font-weight:bold;}", dipLabelFontSize).
    Textf("text.chipSubLabel {font-size:%dpx;}", dipSubLabelFontSize).
    End().
    G().Class("chip").
    // DIP visualisation
    G().Class("chipCase").
    Rect().X(vWidth2-dipWidth2).Y(10).End().
    Circle().CX(vWidth2).CY(10).R(10).End().
    End(). // G.chipCase
    // DIP Labels
    G().Class("chipLabel").
    Attr("transform", "translate(%d %d) rotate(90)", vWidth2, height2).
    SvgText().Y(-7).Class("chipLabel").Text(d.Label).End().
    SvgText().Y(30).Class("chipSubLabel").Text(d.SubLabel).End().
    End(). // G.chipLabel
      // End DIP visualisation
      // left side counting top to bottom
      Sequence(1, d.PinCount, func(pin int, e *html.Element) *html.Element {
        var x, y int
        z := 8 + dipPinWidth
        side := pin <= pinCount2
        if side {
          x = vWidth2 - dipWidth2 - dipPinWidth
          y = (dipPinSpacingV * pin) - dipPinSpacingV2
          z = dipPinWidth - z
        } else {
          x = vWidth2 + dipWidth2
          y = dipPinSpacingV*(d.PinCount-pin+1) - dipPinSpacingV2
        }
        e = dipPin(pin, x, y+8, e)
        e = dipPinLabel(side, x+z, y+25, d.Pins[pin], e)
        return e
      }).
    End(). // G()
    End() // svg

  // =================================================================
  // Convert to FileBuilder & generate the shortcode html file
  return b.
    FileBuilder().
    FileHandler().
    Write(d.Path("content/chipref/reference")+".svg", d.FileInfo.ModTime())
}

func dipPin(pin, x, y int, e *html.Element) *html.Element {
  return e.Rect().
    Class("chipPin").
    X(x).Y(y).
    //Width(dipPinWidth).Height(dipPinHeight).
    //Fill(html.LightGrey).
    //Stroke(html.BLACK).
    End(). // Rect
    SvgText().
    Class("chipPin").
    X(x+dipPinWidth2).Y(y+dipPinHeight2+6).
    //AttrInt("font-size", dipPinFontSize).
    //Attr("text-anchor", "middle").
    Textf("%d", pin).
    End() // SvgText
}

func dipPinLabel(rightAlign bool, x, y int, label string, e *html.Element) *html.Element {
  b := e.G().Attr("transform", "translate(%d %d)", x, y)
  _ = Parse(label, dipPinFontSize, rightAlign, b)
  return b.End()
}

func dipPinLabelV(rightAlign bool, x, y int, label string, e *html.Element) *html.Element {
  b := e.G().Attr("transform", "translate(%d %d) rotate(90deg)", x, y)
  _ = Parse(label, dipPinFontSize, rightAlign, b)
  return b.End()
}
