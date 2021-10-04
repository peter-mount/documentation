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
    Text("g.chipCase {fill:none;stroke:black;stroke-width:2px;}").
    Textf("g.chipCase rect {width:%dpx;height:%dpx;}", dipWidth, height).
    Textf("rect.chipPin {width:%dpx;height:%dpx;fill:%s;stroke:%s;}", 10, dipPinHeight, html.LightGrey, html.BLACK).
    Textf("text.chipPin {font-size:%dpx;text-anchor:middle;}", dipPinFontSize).
    Text("g.chipLabel {fill:lightgrey;text-anchor:middle;}").
    Textf("text.chipLabel {font-size:%dpx;font-weight:bold;}", dipLabelFontSize).
    Textf("text.chipSubLabel {font-size:%dpx;}", dipSubLabelFontSize).
    End().
    G().Class("chip").
    // DIP visualisation
    ClipPath().Attr("id", "chipClip").Rect().X(vWidth2-dipWidth2-1).Y(9).Width(dipWidth+2).Height(height+2).End().End().
    G().Class("chipCase").Attr("clip-path", "url(#chipClip)").
    Rect().X(vWidth2-dipWidth2).Y(10).End().
    Circle().CX(vWidth2).CY(10).R(10).End().
    End(). // G.chipCase
    // DIP Labels
    G().Class("chipLabel").
    Attr("transform", "translate(%d %d) rotate(90)", vWidth2, height2).
    SvgText().Y(-3).Class("chipLabel").Text(d.Label).End().
    SvgText().Y(30).Class("chipSubLabel").Text(d.SubLabel).End().
    End(). // G.chipLabel
      // End DIP visualisation
      // left side counting top to bottom
      Sequence(1, d.PinCount, func(pin int, e *html.Element) *html.Element {
        var x, y int
        side := pin <= pinCount2
        if side {
          x = vWidth2 - dipWidth2 - 9
          y = (dipPinSpacingV * pin) - dipPinSpacingV2
        } else {
          x = vWidth2 + dipWidth2 - dipPinWidth + 10
          y = dipPinSpacingV*(d.PinCount-pin+1) - dipPinSpacingV2
        }
        e = dipPinLabel2(side, pin, x, y+8, 0, d.Pins[pin], e)
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
