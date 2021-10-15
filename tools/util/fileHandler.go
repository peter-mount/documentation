package util

import (
  "bytes"
  "io"
  "io/ioutil"
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

// Bytes returns the content as a []byte
func (a FileHandler) Bytes() ([]byte, error) {
  var buf bytes.Buffer
  err := a(&buf)
  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

// Write writes the file. If the file exists & has the same fileTime then it's content is checked before writing it.
func (a FileHandler) Write(fileName string, fileTime time.Time) error {
  // Do we write or ignore
  writeNow := true

  // File length if it exists
  var fl int64

  if fi, err := os.Stat(fileName); err != nil {
    // If error is not-existing then fall through with writeNow is true
    if !os.IsNotExist(err) {
      return err
    }

    writeNow = true
  } else {
    fl = fi.Size()
    writeNow = fileTime.After(fi.ModTime())
  }

  if !writeNow && fl > 0 {
    // Get our new content as a byte array
    bAry, err := a.Bytes()
    if err != nil {
      return err
    }

    // If the sizes match then compare them for equality
    if int64(len(bAry)) == fl {
      // FIXME do this in a better way?
      // Read the file so compare it as same size
      fBuf, err := ioutil.ReadFile(fileName)
      if err != nil {
        return err
      }

      // If identical then do nothing
      if bytes.Equal(bAry, fBuf) {
        return nil
      }
    }

    // As we have the content use it & write the file
    return ByteFileHandler(bAry).WriteAlways(fileName, fileTime)
  }

  return a.WriteAlways(fileName, fileTime)
}

const (
  referenceDir = "reference"
)

// WriteAlways writes the file regardless of the existing files status
func (a FileHandler) WriteAlways(fileName string, fileTime time.Time) error {
  err := a.writeAlways(fileName, fileTime)
  if err != nil {
    return err
  }

  return GenerateReferenceIndex(fileName, fileTime)
}

func (a FileHandler) writeAlways(fileName string, fileTime time.Time) error {
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

  if fileTime.IsZero() {
    fileTime = time.Now()
  }
  return os.Chtimes(fileName, fileTime, fileTime)
}
