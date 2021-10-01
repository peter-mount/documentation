package chip

import (
  "github.com/peter-mount/documentation/tools/util/html"
)

// pin Generates the list entry for each pin on the dip package
func (d *Definition) pin(i int, e *html.Element) *html.Element {
  return e.LI().
    Span().End(). // The actual pin, number added by CSS
    ChipExpansion(d.Pins[i]). // Chip label expression
    End()
}

// dip Dual Inline Pin chip layout
func dip(d *Definition) error {
  b := html.HtmlBuilder()

  if d.PinCount < 40 {
    b.
      Div().Class("chip").Class("dip%d", d.PinCount).
      Div().Class("position-relative").
      // DIP visualisation
      Div().Class("chip-dip").
      Div().Class("chip-orient").Text(".").End(). // Orientation notch
      Div().Class("chip-label").Text(d.Label).End(). // Main label, usually Manufacturer name
      Div().Class("chip-subLabel").Text(d.SubLabel).End(). // Sub label, usually chip ident
      End(). // end chip-dip
      // The pins
      OL().
      Sequence(1, d.PinCount/2, d.pin). // left side counting top to bottom
      Sequence(d.PinCount, (d.PinCount/2)+1, d.pin). // right side, counting bottom up so reversed in html
      End(). // end OL
      // Convert to FileBuilder
      End(). // end position-relative
        // Optional chip title text underneath the chip visualisation
        If(d.Title != "", func(e *html.Element) *html.Element {
          return e.Div().Class("chip-title").Text(d.Title).End()
        }).
      End() // end chip

  } else {
    // =================================================================
    height := (dipPinSpacingV * d.PinCount / 2) + 20
    height2 := height / 2
    vWidth := 200 + dipWidth
    vWidth2 := vWidth / 2
    b.
      Svg().ViewBox(0, 0, vWidth, height+50).
      G().
      Attr("font-family", html.TextFont).
      // DIP visualisation
      Rect().X(vWidth2-dipWidth2).Y(10).Width(dipWidth).Height(height).Fill("#333").End().
      // Orientation Notch
      Circle().CX(vWidth2).CY(10).R(10).Fill(html.WHITE).End().
      // DIP Labels
      G().
      Fill(html.LightGrey).
      Attr("text-anchor", "middle").
      Attr("transform", "translate(%d %d) rotate(90)", vWidth2, height2).
      SvgText().Y(-7).AttrInt("font-size", dipLabelFontSize).Text(d.Label).Attr("font-weight", "bold").End().
      SvgText().Y(30).AttrInt("font-size", dipSubLabelFontSize).Text(d.SubLabel).End().
      End(). // G() dip label
        // End DIP visualisation
        // left side counting top to bottom
        Sequence(1, d.PinCount/2, func(pin int, e *html.Element) *html.Element {
          return dipPin(pin, vWidth2-dipWidth2-dipPinWidth, (dipPinSpacingV*pin)-dipPinSpacingV2+8, e)
        }).
        Sequence(d.PinCount, (d.PinCount/2)+1, func(pin int, e *html.Element) *html.Element {
          return dipPin(pin, vWidth2+dipWidth2, dipPinSpacingV*(d.PinCount-pin+1)-dipPinSpacingV2+8, e)
        }).
      End(). // G()
      End() // svg
  }

  // =================================================================
  // Convert to FileBuilder & generate the shortcode html file
  return b.
    FileBuilder().
    FileHandler().
    Write(d.Path(), d.FileInfo.ModTime())
}

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

func dipPin(pin, x, y int, e *html.Element) *html.Element {
  return e.Rect().
    X(x).Y(y).
    Width(dipPinWidth).Height(dipPinHeight).
    Fill(html.LightGrey).
    Stroke(html.BLACK).
    End().
    SvgText().
    X(x+dipPinWidth2).Y(y+dipPinHeight2+6).
    AttrInt("font-size", dipPinFontSize).
    Attr("text-anchor", "middle").Textf("%d", pin).End()
}
