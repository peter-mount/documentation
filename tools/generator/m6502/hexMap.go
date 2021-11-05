package m6502

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/html"
  "log"
  "sort"
  "strconv"
  "strings"
)

type HexGrid struct {
  m map[string]*HexMap
}

type HexMap struct {
  Data   [][]*HexCell
  Parent *HexMap // Link to parant Map
}

type HexCell struct {
  Index      string  // Full index value, e.g. 0xEA for NOP
  Row        int     // Row in table, e.g. higher nibble or 0xE for NOP
  Col        int     // Col in table, e.g. lower nibble or 0xA for NOP
  Label      string  // Label to show, e.g. "NOP"
  Addressing string  // Addressing id
  Link       string  // Optional link to a page from this cell
  Size       int     // Size in bytes
  Cycles     int     // Cycle count
  Colour     string  // Optional colour information
  Extension  bool    // if true then this is a prefix to another HexMap
  Parent     *HexMap // Link to parant Map
}

func (c *HexCell) fromOpcode(o *Opcode) {
  c.Label = o.Op
  c.Index = o.Code
  c.Addressing = o.Addressing
  c.Size = o.Bytes.Int()
  c.Cycles = o.Cycles.Int()
  c.Colour = o.Colour
}

// fromOpcodeIgnoreExtension same as fromOpcode except Do not overwrite Extension flag
func (c *HexCell) fromOpcodeIgnoreExtension(o *Opcode) {
  if !c.Extension {
    c.fromOpcode(o)
  }
}

func NewHexGrid() *HexGrid {
  return &HexGrid{m: make(map[string]*HexMap)}
}

func (hg *HexGrid) resolve(code string) *HexMap {

  prefix := ""
  // Ensure code is linked if it's multi-byte (e.g. Z80)
  // This will handle multiple levels like 0xDDCB by op codes due to it recursing through hg.Cell().
  if len(code) > 2 {
    // Mark Cell as an extension
    prefix = code[:len(code)-2]

    // Why Zilog did this but the Z80 has an oddity where we have the operand inside the opcode.
    // For example "RLC (IX+d)" is DDCBnn06 where DDCB is the prefix but the RLC command is 06 AFTER
    // the byte holding +d. Really odd as other IX codes like "LD H,(IX+d)" is DD66nn
    // So to handle this if we see the nn then we fudge it so the HexMap is for the part before nn
    // & contains the codes after the nn.
    lookupPrefix := prefix
    if strings.HasSuffix(prefix, "nn") {
      lookupPrefix = lookupPrefix[:len(lookupPrefix)-2]
    }

    c := hg.Cell(lookupPrefix)
    c.Extension = true
    c.Label = "Instruction Prefix" //fmt.Sprintf("Prefix 0x%s", prefix)
    c.Index = prefix

    prefix = strings.ReplaceAll(prefix, "nn", "")
  }

  if m, exists := hg.m[prefix]; exists {
    return m
  }

  g := &HexMap{}
  for i := 0; i < 16; i++ {
    var r []*HexCell
    for j := 0; j < 16; j++ {
      r = append(r, &HexCell{
        Index:  fmt.Sprintf("%s%X%X", prefix, i, j),
        Row:    i,
        Col:    j,
        Parent: g,
      })
    }
    g.Data = append(g.Data, r)
  }

  log.Printf("*** Prefix %q", prefix)
  g.Parent = hg.m[prefix]
  hg.m[prefix] = g

  return g
}

// Cell resolves the appropriate HexCell for the opCode
func (hg *HexGrid) Cell(code string) *HexCell {
  return hg.resolve(code).cell(code)
}

func (hg *HexGrid) Opcode(a ...*Opcode) *HexGrid {
  for _, o := range a {
    if o != nil {
      /*      i := strings.Index(o.Code, "nn")
              if i > -1 {
                // Special case where Z80 has FDCBnnCC so we want to locate FDCB with
                // prefix FDCBnn
                log.Println("**** ", o)
                c := o.Code[:i]
                log.Println(c)
                m := hg.resolve(c)
                log.Println(m)
                log.Println(m.cell(c))
                hg.Cell(c).
                  fromOpcodeIgnoreExtension(o)
              } else {
      */hg.Cell(o.Code).
        fromOpcodeIgnoreExtension(o)
      /*      }
       */}
  }
  return hg
}

func (hg *HexGrid) FileBuilder() util.FileBuilder {
  return func(slice util.StringSlice) (util.StringSlice, error) {
    slice = append(slice, "hexGrid:")

    var keys []string
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

func fixN(i int) string {
  return fix(fmt.Sprintf("%X", i))
}

func fix(s string) string {
  if len(s)%2 == 1 {
    return "0" + s
  }
  return s
}

// Cell resolves the appropriate HexCell for the opCode
func (g *HexMap) cell(code string) *HexCell {
  s := strings.ReplaceAll(code, "nn", "")
  if len(s) > 2 {
    s = s[len(s)-2:]
  }

  i, err := strconv.ParseInt(s, 16, 64)
  if err != nil {
    log.Printf("*** Failed HexMap.cell(%q): s=%q\n%v", code, s, err)
    return nil
  }

  c := g.Data[(i&0xf0)>>4][i&0x0f]
  c.Index = code

  return c
}

func (g *HexMap) opcode(a ...*Opcode) {
  for _, o := range a {
    if o != nil {
      g.cell(o.Code).
        fromOpcode(o)
    }
  }
}

func (g *HexMap) write(slice util.StringSlice) (util.StringSlice, error) {
  for _, r := range g.Data {
    slice = append(slice, "    -")
    for _, c := range r {
      slice = append(slice, fmt.Sprintf("      - label: %q", c.Label))
      if c.Label != "" {
        slice = append(slice, fmt.Sprintf("        op: %q", c.Index))
      }
      if c.Extension {
        slice = append(slice, fmt.Sprintf("        colour: %q", "grey"))
      } else {
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
        if c.Colour != "" {
          slice = append(slice, fmt.Sprintf("        colour: %q", c.Colour))
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
