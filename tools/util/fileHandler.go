package util

import (
  "io"
  "log"
  "os"
  "path"
  "time"
)

type FileHandler func(w io.Writer) error

func (a FileHandler) Then(b FileHandler) FileHandler {
  return func(w io.Writer) error {
    if err := a(w); err != nil {
      return err
    }
    return b(w)
  }
}

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

func ByteFileHandler(b []byte) FileHandler {
  return func(w io.Writer) error {
    _, err := w.Write(b)
    return err
  }
}

func StringFileHandler(s string) FileHandler {
  return ByteFileHandler([]byte(s))
}

func WriteFile(fileName string, fileTime time.Time, h FileHandler) error {
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

  err = h(f)
  if err != nil {
    return err
  }

  return os.Chtimes(fileName, fileTime, fileTime)
}
