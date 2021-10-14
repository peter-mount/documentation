package m6502

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/html"
)

type HexMap struct {
  Data [][]*HexCell
}

type HexCell struct {
  Index      int    // Full index value, e.g. 0xEA for NOP
  Row        int    // Row in table, e.g. higher nibble or 0xE for NOP
  Col        int    // Col in table, e.g. lower nibble or 0xA for NOP
  Label      string // Label to show, e.g. "NOP"
  Addressing string // Addressing id
  Link       string // Optional link to a page from this cell
  Size       int    // Size in bytes
  Cycles     int    // Cycle count
}

func NewHexGrid() *HexMap {
  g := &HexMap{}
  for i := 0; i < 16; i++ {
    var r []*HexCell
    for j := 0; j < 16; j++ {
      r = append(r, &HexCell{
        Index: (i << 4) | j,
        Row:   i,
        Col:   j,
      })
    }
    g.Data = append(g.Data, r)
  }
  return g
}

func (g *HexMap) Cell(i int) *HexCell {
  return g.Data[(i&0xf0)>>4][i&0x0f]
}

func (g *HexMap) Label(i int, s string) *HexMap {
  g.Cell(i).Label = s
  return g
}

func (g *HexMap) Link(i int, s string) *HexMap {
  g.Cell(i).Link = s
  return g
}

func (g *HexMap) Opcode(a ...*Opcode) *HexMap {
  for _, o := range a {
    if o != nil {
      c := g.Cell(o.Order)
      c.Label = o.Op
      c.Addressing = o.Addressing
      c.Size = o.Bytes.Int()
      c.Cycles = o.Cycles.Int()
    }
  }
  return g
}

func (g *HexMap) FileBuilder() util.FileBuilder {
  return func(slice util.StringSlice) (util.StringSlice, error) {
    slice = append(slice, "hexGrid:")
    for _, r := range g.Data {
      slice = append(slice, "  -")
      for _, c := range r {
        slice = append(slice, fmt.Sprintf("    - label: %q", c.Label))
        if c.Addressing != "" {
          slice = append(slice, fmt.Sprintf("      addressing: %q", c.Addressing))
        }
        if c.Link != "" {
          slice = append(slice, fmt.Sprintf("      link: %q", c.Link))
        }
        if c.Size > 0 {
          slice = append(slice, fmt.Sprintf("      size: %d", c.Size))
        }
        if c.Cycles > 0 {
          slice = append(slice, fmt.Sprintf("      cycles: %d", c.Cycles))
        }
      }
    }
    return slice, nil
  }
}

func (g *HexMap) colSequenceHandler(cellId int, e *html.Element) *html.Element {
  e = e.TH()
  if cellId >= 0 {
    e = e.Textf("%x", cellId)
  }
  return e.End()
}

func (g *HexMap) rowSequenceHandler(rowId int, e *html.Element) *html.Element {
  return e.TR().
    TH().Textf("%x", rowId).End().
      Sequence(0, 15, func(cellId int, e *html.Element) *html.Element {
        cell := g.Data[rowId][cellId]
        e = e.TD()
        if cell.Label == "" {
          e.Class("empty")
        }

        if cell.Link != "" {
          e = e.A().Attr("href", cell.Link)
        }

        e = e.Text(cell.Label)

        if cell.Link != "" {
          e = e.End()
        }

        return e.End()
      }).
    End()
}
