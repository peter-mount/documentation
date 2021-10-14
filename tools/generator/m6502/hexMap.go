package m6502

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/html"
  "sort"
)

type HexGrid struct {
  m map[int]*HexMap
}

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
  Extension  bool   // if true then this is a prefix to another HexMap
}

func NewHexGrid() *HexGrid {
  return &HexGrid{m: make(map[int]*HexMap)}
}

func (hg *HexGrid) resolve(opCode int) *HexMap {

  // Ensure code is linked if it's multi-byte (e.g. Z80)
  // This will handle multiple levels like 3 by op codes due to it recursing
  if opCode >= 0x100 {
    // Mark Cell as an extension
    c := hg.Cell(opCode >> 8)
    if !c.Extension {
      c.Extension = true
      c.Label = fmt.Sprintf("Prefix 0x%X", opCode>>8)
    }
  }

  if m, exists := hg.m[opCode>>8]; exists {
    return m
  }

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

  hg.m[opCode>>8] = g
  return g
}

// Cell resolves the appropriate HexCell for the opCode
func (hg *HexGrid) Cell(i int) *HexCell {
  return hg.resolve(i).cell(i)
}

func (hg *HexGrid) Opcode(a ...*Opcode) *HexGrid {
  for _, o := range a {
    if o != nil {
      c := hg.Cell(o.Order)
      // Do not overwrite Extension flag
      if !c.Extension {
        c.Label = o.Op
        c.Addressing = o.Addressing
        c.Size = o.Bytes.Int()
        c.Cycles = o.Cycles.Int()
      }
    }
  }
  return hg
}

func (hg *HexGrid) FileBuilder() util.FileBuilder {
  return func(slice util.StringSlice) (util.StringSlice, error) {
    slice = append(slice, "hexGrid:")

    var keys []int
    for k, _ := range hg.m {
      keys = append(keys, k)
    }
    sort.SliceStable(keys, func(i, j int) bool {
      return keys[i] < keys[j]
    })

    var err error
    for _, key := range keys {
      slice = append(slice, fmt.Sprintf("  %q:", fix(key)))
      slice, err = hg.m[key].write(slice)
      if err != nil {
        return nil, err
      }
    }

    return slice, nil
  }
}

func fix(i int) string {
  s := fmt.Sprintf("%X", i)
  if len(s)%2 == 1 {
    return "0" + s
  }
  return s
}

// Cell resolves the appropriate HexCell for the opCode
func (g *HexMap) cell(i int) *HexCell {
  c := g.Data[(i&0xf0)>>4][i&0x0f]
  if c.Index != i {
    c.Index = i
  }
  return c
}

func (g *HexMap) opcode(a ...*Opcode) {
  for _, o := range a {
    if o != nil {
      c := g.cell(o.Order)
      c.Label = o.Op
      c.Addressing = o.Addressing
      c.Size = o.Bytes.Int()
      c.Cycles = o.Cycles.Int()
    }
  }
}

func (g *HexMap) write(slice util.StringSlice) (util.StringSlice, error) {
  for _, r := range g.Data {
    slice = append(slice, "    -")
    for _, c := range r {
      slice = append(slice, fmt.Sprintf("      - label: %q", c.Label))
      if c.Label != "" {
        slice = append(slice, fmt.Sprintf("        op: %q", fix(c.Index)))
      }
      if !c.Extension {
        if c.Addressing != "" {
          slice = append(slice, fmt.Sprintf("        addressing: %q", c.Addressing))
        }
        if c.Link != "" {
          slice = append(slice, fmt.Sprintf("        link: %q", c.Link))
        }
        if c.Size > 0 {
          slice = append(slice, fmt.Sprintf("        size: %d", c.Size))
        }
        if c.Cycles > 0 {
          slice = append(slice, fmt.Sprintf("        cycles: %d", c.Cycles))
        }
      }
    }
  }
  return slice, nil
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
