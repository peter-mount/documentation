package util

import (
  "io"
  "log"
  "os"
  "path"
  "time"
)

type FileHandler func(w io.Writer) error

// Then returns a new FileHandler which will run the left hand side first then the right hand side.
func (a FileHandler) Then(b FileHandler) FileHandler {
  return func(w io.Writer) error {
    if err := a(w); err != nil {
      return err
    }
    return b(w)
  }
}

// FileHandlerOf returns a FileHandler which will run each supplied FileHandler in sequence
func FileHandlerOf(h ...FileHandler) FileHandler {
  switch len(h) {
  case 0:
    return func(_ io.Writer) error {
      return nil
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

// ByteFileHandler returns a FileHandler which will write a []byte
func ByteFileHandler(b []byte) FileHandler {
  return func(w io.Writer) error {
    _, err := w.Write(b)
    return err
  }
}

// StringFileHandler returns a FileHandler which will write a string
func StringFileHandler(s string) FileHandler {
  return ByteFileHandler([]byte(s))
}

func (a FileHandler) Write(fileName string, fileTime time.Time) error {
  log.Printf("Writing %s", fileName)
  err := os.MkdirAll(path.Dir(fileName), 0755)
  if err != nil {
    return err
  }

  f, err := os.Create(fileName)
  if err != nil {
    return err
  }
  defer f.Close()

  err = a(f)
  if err != nil {
    return err
  }

  return os.Chtimes(fileName, fileTime, fileTime)
}
