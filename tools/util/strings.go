package util

import (
  "io"
  "sort"
  "strings"
)

type StringHandler func(string) error

// WithStringHandler starts a new StringHandler
func WithStringHandler() StringHandler {
  return func(_ string) error {
    return nil
  }
}

// Then joins two string handlers so the left-hand side runs first then the right-hand side.
func (a StringHandler) Then(b StringHandler) StringHandler {
  return func(s string) error {
    err := a(s)
    if err != nil {
      return err
    }
    return b(s)
  }
}

// Do invokes a StringHandler with a specific string
func (a StringHandler) Do(s string) error {
  return a(s)
}

type StringSlice []string

// StringSliceOf returns a StringSlice from a slice of strings
func StringSliceOf(s ...string) StringSlice {
  return s
}

// ForEach invokes a StringHandler for each entry in a StringSlice. First one to return an error terminates the loop.
func (s StringSlice) ForEach(f StringHandler) error {
  for _, b := range s {
    err := f(b)
    if err != nil {
      return err
    }
  }
  return nil
}

// FileHandler returns a FileHandler which will contain the StringSlice with each entry being a single line of the output.
func (s StringSlice) FileHandler() FileHandler {
  return StringFileHandler(strings.Join(s, "\n"))
}

// Write will write a StringSlice to a writer with each entry being a single line.
func (s StringSlice) Write(w io.Writer) error {
  _, err := w.Write([]byte(strings.Join(s, "\n")))
  return err
}

// IsEmpty returns true if a StringSlice has no entries
func (s StringSlice) IsEmpty() bool {
  return len(s) == 0
}

// Join returns the content of a StringSlice with the specified separator between each entry
func (s StringSlice) Join(sep string) string {
  return strings.Join(s, sep)
}

// Join2 is similar to Join except it adds the prefix & suffix to the final result
func (s StringSlice) Join2(prefix, suffix, sep string) string {
  return prefix + s.Join(sep) + suffix
}

// Sort sorts the string slice using a case-insensitive comparator.
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
