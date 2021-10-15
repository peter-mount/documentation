package autodoc

import (
  "github.com/peter-mount/documentation/tools/util"
  "io"
  "log"
  "os"
  "path"
  "strings"
  "time"
)

type Builder interface {
  Comment(string, ...interface{}) Builder
  Header(string, string, string) Builder
  Newline() Builder
  Separator() Builder
  InvokeTopic(string, TopicHandler) Builder
  Invoke(Handler) Builder
  Using(p Provider) Builder
  Do() error
}

type Provider func(string, string, time.Time) Builder

type Handler func(Builder) error

type TopicHandler func(string, string) Handler

func InitBuilder(fileName string, modified time.Time) (io.WriteCloser, error) {
  writeNow := true

  if fi, err := os.Stat(fileName); err != nil {
    // If error is not-existing then fall through with writeNow is true
    if !os.IsNotExist(err) {
      return nil, err
    }

    writeNow = true
  } else {
    writeNow = modified.After(fi.ModTime())
  }

  if err := os.MkdirAll(path.Dir(fileName), 0755); err != nil {
    return nil, err
  }

  if err := util.GenerateReferenceIndices(fileName, modified); err != nil {
    return nil, err
  }

  if writeNow || true {
    log.Println("Creating", fileName)
    f, err := os.Create(fileName)
    if err != nil {
      return nil, err
    }

    return f, nil
  }

  return nil, nil
}

func CloseBuilder(err error, w io.WriteCloser, fileName string, modified time.Time) error {
  if w != nil {
    log.Println("Finishing", fileName)
    if err == nil {
      err = w.Close()
    }

    if err == nil {
      if modified.IsZero() {
        modified = time.Now()
      }

      err = os.Chtimes(fileName, modified, modified)
    }
  }

  return err
}

func Write(err *error, w *io.WriteCloser, s string) {
  if *w != nil && *err == nil {
    if !strings.HasSuffix(s, "\n") {
      s = s + "\n"
    }
    _, *err = (*w).Write([]byte(s))
  }
}

type unionBuilder struct {
  dir      string
  file     string
  modified time.Time
  src      []Builder
}

func For(dir, file string, modified time.Time) Builder {
  return &unionBuilder{dir: dir, file: file, modified: modified}
}

func (u *unionBuilder) Using(p Provider) Builder {
  u.src = append(u.src, p(u.dir, u.file, u.modified))
  return u
}

func (u *unionBuilder) Invoke(handler Handler) Builder {
  for _, b := range u.src {
    b.Invoke(handler)
  }
  return u
}

func (u *unionBuilder) InvokeTopic(t string, h TopicHandler) Builder {
  for _, b := range u.src {
    b.InvokeTopic(t, h)
  }
  return u
}

func (u *unionBuilder) Do() error {
  for _, b := range u.src {
    err := b.Do()
    if err != nil {
      return err
    }
  }
  return nil
}

func (u *unionBuilder) Comment(s string, a ...interface{}) Builder {
  for _, b := range u.src {
    b.Comment(s, a...)
  }
  return u
}

func (u *unionBuilder) Header(s string, s2 string, s3 string) Builder {
  for _, b := range u.src {
    b.Header(s, s2, s3)
  }
  return u
}

func (u *unionBuilder) Separator() Builder {
  for _, b := range u.src {
    b.Separator()
  }
  return u
}

func (u *unionBuilder) Newline() Builder {
  for _, b := range u.src {
    b.Newline()
  }
  return u
}
