package bbc

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "sort"
)

type Osbyte struct {
  call   int                         // Call 0..255
  params map[interface{}]interface{} // Params
}

func (b *BBC) extractOsbyte(osbyte interface{}) {
  if a, ok := osbyte.([]interface{}); ok {
    for _, e := range a {
      if m, ok := e.(map[interface{}]interface{}); ok {
        if v, ok := util.DecodeInt(m["int"], 0); ok {
          o := &Osbyte{
            call:   v,
            params: m,
          }
          b.osbyte = append(b.osbyte, o)
        }
      }
    }
  }
}

func (b *BBC) writeOsbyteIndex(book *hugo.Book) error {
  sort.SliceStable(b.osbyte, func(i, j int) bool {
    return b.osbyte[i].call < b.osbyte[j].call
  })

  r := Output{Nometa: true}
  for _, o := range b.osbyte {
    r.Osbyte = append(r.Osbyte, o.params)
  }

  return util.WriteReferenceYaml(
    book.ContentPath(),
    "osbyte",
    "OSByte calls",
    "OSByte &FFF4 calls",
    10,
    book.Modified(),
    r)
}
