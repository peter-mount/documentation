package util

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util/strings"
  "gopkg.in/yaml.v2"
  "io"
  "path"
  "strconv"
  "time"
)

// FileBuilder creates text files
type FileBuilder func(strings.StringSlice) (strings.StringSlice, error)

// FileBuilderOf returns a FileBuilder which consists of each provided FileBuilder in sequence
func FileBuilderOf(h ...FileBuilder) FileBuilder {
  switch len(h) {
  case 0:
    return func(slice strings.StringSlice) (strings.StringSlice, error) {
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
  return func(s strings.StringSlice) (strings.StringSlice, error) {
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
  return func(slice strings.StringSlice) (strings.StringSlice, error) {
    slice = append(slice, "---")
    slice, err := a(slice)
    if err != nil {
      return nil, err
    }
    slice = append(slice, "---")
    return slice, nil
  }
}

func (a FileBuilder) Append(s string) FileBuilder {
  return a.Then(func(slice strings.StringSlice) (strings.StringSlice, error) {
    return append(slice, s), nil
  })
}

func (a FileBuilder) Appendf(s string, args ...interface{}) FileBuilder {
  return a.Then(func(slice strings.StringSlice) (strings.StringSlice, error) {
    return append(slice, fmt.Sprintf(s, args...)), nil
  })
}

// Yaml appends the supplied value as YAML
func (a FileBuilder) Yaml(val interface{}) FileBuilder {
  return func(slice strings.StringSlice) (strings.StringSlice, error) {
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
func ReferenceFileBuilder(title, desc, layout string, weight int, t time.Time) FileBuilder {
  return func(ary strings.StringSlice) (strings.StringSlice, error) {
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
    if !t.IsZero() {
      ary = append(ary, fmt.Sprintf("lastmod: %q", t.Format(time.RFC3339)))
    }
    return ary, nil
  }
}

func BlankFileBuilder() FileBuilder {
  return func(slice strings.StringSlice) (strings.StringSlice, error) {
    return slice, nil
  }
}
