package util

import (
  "io"
  "strings"
)

type StringHandler func(string) error

func (a StringHandler) Then(b StringHandler) StringHandler {
  return func(s string) error {
    err := a(s)
    if err != nil {
      return err
    }
    return b(s)
  }
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
