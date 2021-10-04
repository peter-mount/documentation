package chip

import (
  "github.com/peter-mount/documentation/tools/util/html"
  "strings"
)

func lccc(d *Definition) error {
  // For now lccc is just an alias for cc
  return cc(d)
}

func plcc(d *Definition) error {
  // For now plcc is just an alias for cc
  return cc(d)
}

func qfp(d *Definition) error {
  // For now qfp is just an alias for cc
  return cc(d)
}

func cc(d *Definition) error {
  b := html.HtmlBuilder()

  pinCount4 := d.PinCount / 4

  height := dipPinSpacingV * (pinCount4 + 3)
  height2 := height / 2
  width := height // These are square chips
  width2 := width / 2

  vWidth := 200 + width
  vWidth2 := vWidth / 2
  vHeight := vWidth
  vHeight2 := vHeight / 2

  b = b.
    Svg().ViewBox(0, 0, vWidth, vHeight).Width(vWidth).
    Style().Attr("type", "text/css").
    // TODO check text-decoration-thickness as seems broken in Chrome 94 but works in Firefox 92
    Text("tspan.not {text-decoration: overline;webkit-text-decoration-thickness: 2px;text-decoration-thickness: 2px;}").
    Textf(".chip {font-family:%s;}", html.TextFont).
    Textf("g.chipCase polygon {stroke:black;fill:none;stroke-width: 2px;}").
    Text(".chipCase circle {fill:white;}").
    Textf("rect.chipPin {width:%dpx;height:%dpx;fill:%s;stroke:%s;stroke-dasharray:3 3;}", 10, dipPinHeight, "#eee", html.BLACK).
    Textf("text.chipPin {font-size:%dpx;text-anchor:middle;}", dipPinFontSize).
    Textf("g.chipLabel {fill:black;text-anchor:middle;}").
    Textf("text.chipLabel {font-size:%dpx;font-weight:bold;}", dipLabelFontSize).
    Textf("text.chipSubLabel {font-size:%dpx;}", dipSubLabelFontSize).
    Textf("text.chipTypeLabel {font-size:%dpx;}", dipSubLabelFontSize/2).
    End(). // style
    G().Class("chip").
    // DIP Labels
    G().Class("chipLabel").
    Attr("transform", "translate(%d %d)", vWidth2, vHeight2).
    SvgText().Y(-7).Class("chipLabel").Text(d.Label).End().
    SvgText().Y(30).Class("chipSubLabel").Text(d.SubLabel).End().
    SvgText().Y(height2-40).Class("chipTypeLabel").Text(strings.ToUpper(d.Type)).End().
    End(). // G.chipLabel
      Sequence(1, d.PinCount, func(pin int, e *html.Element) *html.Element {
        // Number of pins per side
        pinSizeCount := d.PinCount / 4

        // Actual pin position accounting for pinOffset
        pinIndex := pin + d.PinOffset
        for pinIndex < 1 {
          pinIndex += d.PinCount
        }

        // Position of pin on associated side
        pinSideIndex := (pinIndex - 1) % pinSizeCount // pin # on side

        // side from 0..3 counter-clockwise from left
        pinSide := (pinIndex - 1) / pinSizeCount

        // Convert pinSideIndex to account for correct counter-clockwise pin ordering
        if pinSide > 1 {
          pinSideIndex = pinSizeCount - pinSideIndex - 1
        }

        // Calculate position of pins
        align := pinSide < 2 // Side 0 & 1 are right aligned
        x := vWidth2 - (dipPinWidth2 / 2) + (dipPinWidth * (pinSideIndex - pinSizeCount/2))
        y := vHeight2 - (dipPinSpacingV2 / 2) + (dipPinSpacingV * (pinSideIndex - pinSizeCount/2))
        rotate := 0 // text rotation

        // Adjust
        switch (pinIndex - 1) / pinSizeCount {
        case 0:
          x = vWidth2 - width2 + 1
        case 1:
          y = vHeight2 + height2 - 1
          rotate = 1
        case 2:
          x = vWidth2 + width2 - dipPinWidth - 1
        case 3:
          y = vHeight2 - height2 + dipPinSpacingV + 3
          rotate = 1
        }

        e = dipPinLabel2(align, pin, x, y, rotate, d.Pins[pin], e)
        return e
      }).
    // DIP visualisation - last so it overlays the pins
    G().Class("chipCase").
    Polygon().
    Point(vWidth2-width2+16, vHeight2-height2).
    Point(vWidth2-width2, vHeight2-height2+16).
    Point(vWidth2-width2, vHeight2+height2).
    Point(vWidth2+width2, vHeight2+height2).
    Point(vWidth2+width2, vHeight2-height2).
    End(). //Polygon()
    End() // G.chipCase

  // =================================================================
  // Convert to FileBuilder & generate the shortcode html file
  return b.
    End(). // Svg
    FileBuilder().
    FileHandler().
    Write(d.Path("content/chipref/reference")+".svg", d.FileInfo.ModTime())
}
