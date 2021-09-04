package util

import (
  "gopkg.in/yaml.v2"
  "io"
  "path"
  "strconv"
)

// FileBuilder creates text files
type FileBuilder func(StringSlice) (StringSlice, error)

// FileBuilderOf returns a FileBuilder which consists of each provided FileBuilder in sequence
func FileBuilderOf(h ...FileBuilder) FileBuilder {
  switch len(h) {
  case 0:
    return func(slice StringSlice) (StringSlice, error) {
      return slice, nil
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

// Then returns a FileBuilder which runs the first then second FileBuilder in sequence
func (a FileBuilder) Then(b FileBuilder) FileBuilder {
  return func(s StringSlice) (StringSlice, error) {
    s, err := a(s)
    if err != nil {
      return nil, err
    }
    return b(s)
  }
}

// WrapAsFrontMatter wraps a FileBuilder so that it's generated content is wrapped by a pair of "---" lines
// suitable for Hugo's front matter
func (a FileBuilder) WrapAsFrontMatter() FileBuilder {
  return func(slice StringSlice) (StringSlice, error) {
    slice = append(slice, "---")
    slice, err := a(slice)
    if err != nil {
      return nil, err
    }
    slice = append(slice, "---")
    return slice, nil
  }
}

// Yaml appends the supplied value as YAML
func (a FileBuilder) Yaml(val interface{}) FileBuilder {
  return func(slice StringSlice) (StringSlice, error) {
    slice, err := a(slice)
    if err != nil {
      return nil, err
    }

    c, err := yaml.Marshal(val)
    if err != nil {
      return nil, err
    }

    slice = append(slice, string(c[:]))

    return slice, nil
  }
}

// FileHandler converts the FileBuilder into a FileHandler
func (a FileBuilder) FileHandler() FileHandler {
  return func(w io.Writer) error {
    s, err := a(nil)
    if err != nil {
      return err
    }
    return s.Write(w)
  }
}

// ReferenceFilename returns the file name for a reference file.
func ReferenceFilename(dir, name, fileName string) string {
  if name == "" {
    return path.Join(dir, "reference", fileName)
  }
  return path.Join(dir, "reference", name, fileName)
}

// ReferenceFileBuilder returns a FileBuilder with standard front matter for hugo pages
func ReferenceFileBuilder(title, desc, layout string, weight int) FileBuilder {
  return func(ary StringSlice) (StringSlice, error) {
    ary = append(ary,
      "# This file is generated.",
      "# To edit, change the files under content then run the generator.",
      "type: \""+layout+"\"",
      "title: \""+title+"\"",
      "linkTitle: \""+title+"\"",
      "weight: "+strconv.Itoa(weight),
      "description: \""+desc+"\"",
      "notitle: \"don't show instructions header in table\"",
    )
    return ary, nil
  }
}
