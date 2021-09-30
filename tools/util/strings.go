package util

import (
  "io"
  "sort"
  "strings"
)

type StringHandler func(string) error

func WithStringHandler() StringHandler {
  return func(_ string) error {
    return nil
  }
}

func (a StringHandler) Then(b StringHandler) StringHandler {
  return func(s string) error {
    err := a(s)
    if err != nil {
      return err
    }
    return b(s)
  }
}

func (a StringHandler) Do(s string) error {
  return a(s)
}

type StringSlice []string

func (s StringSlice) ForEach(f StringHandler) error {
  for _, b := range s {
    err := f(b)
    if err != nil {
      return err
    }
  }
  return nil
}

func (s StringSlice) FileHandler() FileHandler {
  return StringFileHandler(strings.Join(s, "\n"))
}

func (s StringSlice) Write(w io.Writer) error {
  _, err := w.Write([]byte(strings.Join(s, "\n")))
  return err
}

func (s StringSlice) IsEmpty() bool {
  return len(s) == 0
}

func (s StringSlice) Join(sep string) string {
  return strings.Join(s, sep)
}

func (s StringSlice) Join2(prefix, suffix, sep string) string {
  return prefix + s.Join(sep) + suffix
}

func (s StringSlice) Sort() StringSlice {
  sort.SliceStable(s, func(i, j int) bool {
    return strings.ToLower(s[i]) < strings.ToLower(s[j])
  })
  return s
}

type StringSliceHandler func(StringSlice) (StringSlice, error)

func (a StringSliceHandler) Then(b StringSliceHandler) StringSliceHandler {
  return func(s StringSlice) (StringSlice, error) {
    s, err := a(s)
    if err != nil {
      return nil, err
    }
    return b(s)
  }
}

func (a StringSliceHandler) FileHandler() FileHandler {
  return func(w io.Writer) error {
    s, err := a(nil)
    if err != nil {
      return err
    }

    return s.Write(w)
  }
}
