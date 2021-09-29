package chip

import (
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/html"
)

func (d *Definition) pin(i int, e *html.Element) {
  e.LI().
    Span().End().
    ChipExpansion(d.Pins[i]).
    End()
}

// dip Dual Inline Pin chip layout
func dip(d *Definition) error {
  return util.BlankFileBuilder().
      Then(html.HtmlBuilder().
        Div().Class("chip").Class("dip%d", d.PinCount).
        Div().Class("position-relative").
        // DIP visualisation
        Div().Class("chip-dip").
        Div().Class("chip-orient").Text(".").End().
        Div().Class("chip-label").Text(d.Label).End().
        Div().Class("chip-subLabel").Text(d.SubLabel).End().
        End().
        // The pins
        OL().
        Sequence(1, d.PinCount/2, d.pin).
        Sequence(d.PinCount, (d.PinCount/2)+1, d.pin).
        End().
        // Convert to FileBuilder
        End().
        End().
        FileBuilder()).
    FileHandler().
    Write(d.Path(), d.FileInfo.ModTime())
}
