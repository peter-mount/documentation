package bbc

import (
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "log"
  "sort"
  "strconv"
  "strings"
)

type Api struct {
  call     int    // Call 0..255
  Name     string `yaml:"name"`
  Addr     string `yaml:"addr"`
  Indirect string `yaml:"indirect,omitempty"`
  Title    string `yaml:"title"`
  params   interface{}
}

func (b *BBC) extractApi(_ *hugo.FrontMatter, api interface{}) error {
  return util.ForEachInterface(api, func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      v := &Api{
        Name:   util.DecodeString(m["name"], ""),
        Addr:   util.DecodeString(m["addr"], ""),
        Title:  util.DecodeString(m["title"], ""),
        params: e,
      }

      if i, err := strconv.ParseInt(v.Addr, 16, 64); err != nil {
        log.Printf("Failed to parse addr \"%s\" for %s", v.Addr, v.Name)
        return err
      } else {
        v.call = int(i)
      }

      b.api = append(b.api, v)

      return nil
    })
  })
}

func (b *BBC) writeAPIIndex(book *hugo.Book) error {
  sort.SliceStable(b.api, func(i, j int) bool {
    return b.api[i].call < b.api[j].call
  })

  r := Output{Nometa: true}
  for _, o := range b.api {
    r.Api = append(r.Api, o.params)
  }

  return util.ReferenceFileBuilder(
    "MOS API by address",
    "MOS API by address",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "api", "_index.html"), book.Modified())
}

func (b *BBC) writeAPINameIndex(book *hugo.Book) error {
  sort.SliceStable(b.api, func(i, j int) bool {
    return strings.ToLower(b.api[i].Name) < strings.ToLower(b.api[j].Name)
  })

  r := Output{Nometa: true}
  for _, o := range b.api {
    r.Api = append(r.Api, o.params)
  }

  return util.ReferenceFileBuilder(
    "MOS API by name",
    "MOS API by name",
    "manual",
    10,
  ).
    Yaml(r).
    WrapAsFrontMatter().
    FileHandler().
    Write(util.ReferenceFilename(book.ContentPath(), "apiName", "_index.html"), book.Modified())
}
