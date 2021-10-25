package bbc

import (
  "context"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
  "log"
  "strconv"
)

type Api struct {
  call     int            // Call 0..255
  Name     string         `yaml:"name"`
  Addr     string         `yaml:"addr"`
  Indirect string         `yaml:"indirect,omitempty"`
  Title    string         `yaml:"title"`
  Entry    FunctionParams `yaml:"entry"`
  Exit     FunctionParams `yaml:"exit"`
  params   interface{}
}

func (b *BBC) extractApi(ctx context.Context, _ *hugo.FrontMatter) error {
  return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
    return util.IfMap(e, func(m map[interface{}]interface{}) error {
      v := &Api{
        Name:     util.DecodeString(m["name"], ""),
        Addr:     util.DecodeString(m["addr"], ""),
        Indirect: util.DecodeString(m["indirect"], ""),
        Title:    util.DecodeString(m["title"], ""),
        params:   e,
      }

      if i, err := strconv.ParseInt(v.Addr, 16, 64); err != nil {
        log.Printf("Failed to parse addr \"%s\" for %s", v.Addr, v.Name)
        return err
      } else {
        v.call = int(i)
      }

      if err := util.IfMapEntry(m, "entry", v.Entry.decode); err != nil {
        return err
      }

      if err := util.IfMapEntry(m, "exit", v.Exit.decode); err != nil {
        return err
      }

      b.api = append(b.api, v)

      return nil
    })
  })
}
