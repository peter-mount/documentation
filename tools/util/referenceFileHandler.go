package util

import (
  "gopkg.in/yaml.v2"
  "io"
  "path"
  "strconv"
  "time"
)

type ReferenceFileHandler func(StringSlice) (StringSlice, error)

func ReferenceFileHandlerOf(h ...ReferenceFileHandler) ReferenceFileHandler {
  switch len(h) {
  case 0:
    return func(slice StringSlice) (StringSlice, error) {
      return nil, nil
    }
  case 1:
    return h[0]
  default:
    a := h[0]
    for _, b := range h[1:] {
      a = a.Then(b)
    }
    return a
  }
}

func (a ReferenceFileHandler) Then(b ReferenceFileHandler) ReferenceFileHandler {
  return func(s StringSlice) (StringSlice, error) {
    s, err := a(s)
    if err != nil {
      return nil, err
    }
    return b(s)
  }
}

func (a ReferenceFileHandler) FileHandler() FileHandler {
  return func(w io.Writer) error {
    s, err := a(nil)
    if err != nil {
      return err
    }
    return s.Write(w)
  }
}

func WriteReferenceFile(dir, name, title, desc, layout string, weight int, fileTime time.Time, handler ReferenceFileHandler) error {
  return WriteFile(
    path.Join(dir, "reference", name, "_index.html"),
    fileTime,
    ReferenceFileHandlerOf(func(a StringSlice) (StringSlice, error) {

      a = append(a,
        "---",
        "# This file is generated.",
        "# To edit, change the files under content then run the generator.",
        "type: \""+layout+"\"",
        "title: \""+title+"\"",
        "linkTitle: \""+title+"\"",
        "weight: "+strconv.Itoa(weight),
        "description: \""+desc+"\"",
        "notitle: \"don't show instructions header in table\"",
      )

      a, err := handler(a)
      if err != nil {
        return nil, err
      }

      a = append(a, "---")

      return a, nil
    }).
      FileHandler())
}

func WriteReferenceYaml(dir, name, title, desc, layout string, weight int, fileTime time.Time, val interface{}) error {
  return WriteReferenceFile(dir, name, title, desc, layout, weight, fileTime, func(a StringSlice) (StringSlice, error) {
    c, err := yaml.Marshal(val)
    if err != nil {
      return nil, err
    }
    a = append(a, string(c[:]))

    return a, nil
  })
}
