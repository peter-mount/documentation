package html

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "strings"
)

type Element struct {
  name     string
  class    util.StringSlice
  parent   *Element
  children []*Element
  text     string
}

func (e *Element) Class(c string, a ...interface{}) *Element {
  e.class = append(e.class, fmt.Sprintf(c, a...))
  return e
}

func (e *Element) String() string {
  var a []string

  if e.text != "" {
    a = append(a, e.text)
  } else {
    if e.name != "" {
      a = append(a, "<", e.name)
      if !e.class.IsEmpty() {
        a = append(a, e.class.Join2(" class=\"", "\"", " "))
      }
      a = append(a, ">")
    }

    for _, child := range e.children {
      a = append(a, child.String())
    }

    if e.name != "" {
      a = append(a, "</", e.name, ">")
    }
  }

  return strings.Join(a, "")
}

func HtmlBuilder() *Element {
  return &Element{}
}

func (e *Element) End() *Element {
  if e != nil && e.parent != nil {
    return e.parent
  }
  return e
}

func (e *Element) RootElement() *Element {
  n := e
  for n.parent != nil {
    n = n.parent
  }
  return n
}

func (e *Element) Textf(f string, s ...interface{}) *Element {
  e.element("").text = fmt.Sprintf(f, s...)
  return e
}

func (e *Element) Text(s ...string) *Element {
  for _, a := range s {
    e.element("").text = a
  }
  return e
}

func (e *Element) FileBuilder() util.FileBuilder {
  n := e.RootElement()
  return func(slice util.StringSlice) (util.StringSlice, error) {
    return append(slice, n.String()), nil
  }
}

func (e *Element) Sequence(start, end int, f func(int, *Element)) *Element {
  if start < end {
    for i := start; i <= end; i++ {
      f(i, e)
    }
  } else {
    for i := start; i >= end; i-- {
      f(i, e)
    }
  }
  return e
}

func (e *Element) Exec(f func(*Element)) *Element {
  f(e)
  return e
}

func (e *Element) element(n string) *Element {
  child := &Element{
    name:   n,
    parent: e,
  }
  e.children = append(e.children, child)
  return child
}

func (e *Element) Div() *Element {
  return e.element("div")
}

func (e *Element) Span() *Element {
  return e.element("span")
}

func (e *Element) OL() *Element {
  return e.element("ol")
}

func (e *Element) LI() *Element {
  return e.element("li")
}

func (e *Element) Sub() *Element {
  return e.element("sub")
}

func (e *Element) Sup() *Element {
  return e.element("sup")
}
