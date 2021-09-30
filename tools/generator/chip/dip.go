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
  return html.HtmlBuilder().
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
    End(). // end chip
    // Convert to FileBuilder & generate the shortcode html file
    FileBuilder().
    FileHandler().
    Write(d.Path(), d.FileInfo.ModTime())
}
