package autodoc

import (
  "fmt"
  "io"
  "path"
  "strings"
  "time"
)

type zAsm struct {
  fileName string
  modified time.Time
  w        io.WriteCloser
  err      error
}

func (b *zAsm) valid() bool {
  return b.err == nil && b.w != nil
}

func ZAsm(dir, file string, modified time.Time) Builder {
  fileName := path.Join(dir, "zasm", file+".z80")
  w, err := InitBuilder(fileName, modified)
  return &zAsm{
    fileName: fileName,
    modified: modified,
    w:        w,
    err:      err,
  }
}

func (b *zAsm) Using(Provider) Builder {
  panic("not implemented")
}

func (b *zAsm) Invoke(handler Handler) Builder {
  if b.err == nil {
    b.err = handler(b)
  }
  return b
}

func (b *zAsm) InvokeTopic(t string, h TopicHandler) Builder {
  return b.Invoke(h(t, b.fileName))
}

func (b *zAsm) Do() error {
  return CloseBuilder(b.err, b.w, b.fileName, b.modified)
}

func (b *zAsm) write(s string) Builder {
  Write(&b.err, &b.w, s)
  return b
}

func (b *zAsm) append(c, f string, a ...interface{}) Builder {
  s := fmt.Sprintf(f, a...)
  if c != "" {
    s = fmt.Sprintf("%-32s ; %s", s, c)
  }
  return b.write(s)
}

func (b *zAsm) Comment(s string, a ...interface{}) Builder {
  // This form of comment starts at the beginning of the line
  return b.write(fmt.Sprintf("; "+s, a...))
}

func (b *zAsm) Header(label, value, comment string) Builder {
  // Value with no label is invalid so ignore
  if label != "" && value != "" {
    return b.append(comment, "%-16s equ %s", label, value)
  }
  return b
}

func (b *zAsm) Newline() Builder {
  return b.write("")
}

func (b *zAsm) Separator() Builder {
  return b.write("; " + strings.Repeat("*", 75))
}
