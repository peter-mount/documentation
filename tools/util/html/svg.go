package html

func (e *Element) Svg() *Element {
  return e.element("svg").
    Attr("xmlns", "http://www.w3.org/2000/svg").
    Attr("version", "1.1")
}

func (e *Element) ViewBox(x, y, width, height int) *Element {
  return e.Attr("viewBox", "%d %d %d %d", x, y, width, height)
}

func (e *Element) R(v int) *Element {
  return e.AttrInt("r", v)
}

func (e *Element) X(v int) *Element {
  return e.AttrInt("x", v)
}

func (e *Element) Y(v int) *Element {
  return e.AttrInt("y", v)
}

func (e *Element) CX(v int) *Element {
  return e.AttrInt("cx", v)
}

func (e *Element) CY(v int) *Element {
  return e.AttrInt("cy", v)
}

func (e *Element) DX(v int) *Element {
  return e.AttrInt("dx", v)
}

func (e *Element) DY(v int) *Element {
  return e.AttrInt("dy", v)
}

func (e *Element) Width(v int) *Element {
  return e.AttrInt("width", v)
}

func (e *Element) Height(v int) *Element {
  return e.AttrInt("height", v)
}

func (e *Element) Fill(v string, a ...interface{}) *Element {
  return e.Attr("fill", v, a...)
}

func (e *Element) Stroke(v string, a ...interface{}) *Element {
  return e.Attr("stroke", v, a...)
}

func (e *Element) StrokeWidth(v string, a ...interface{}) *Element {
  return e.Attr("stroke-width", v, a...)
}

// SvgText is the text element, cannot use Text() as that adds plain text
func (e *Element) SvgText() *Element {
  return e.element("text")
}

func (e *Element) TSpan() *Element {
  return e.element("tspan")
}

func (e *Element) Circle() *Element {
  return e.element("circle")
}

func (e *Element) G() *Element {
  return e.element("g")
}

func (e *Element) Rect() *Element {
  return e.element("rect")
}
