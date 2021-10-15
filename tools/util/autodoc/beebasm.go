package autodoc

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "io"
  "strings"
  "time"
)

type beebAsm struct {
  fileName string
  modified time.Time
  w        io.WriteCloser
  err      error
}

func (b *beebAsm) valid() bool {
  return b.err == nil && b.w != nil
}

func BeebAsm(dir, file string, modified time.Time) Builder {
  fileName, w, err := InitBuilder(dir, file, modified, "beebasm", "asm", "BeebASM", "BeebASM", "Files for the BeebASM assembler")

  return &beebAsm{
    fileName: fileName,
    modified: modified,
    w:        w,
    err:      err,
  }
}

func (b *beebAsm) Using(Provider) Builder {
  panic("not implemented")
}

func (b *beebAsm) Invoke(handler Handler) Builder {
  if b.err == nil {
    b.err = handler(b)
  }
  return b
}

func (b *beebAsm) InvokeTopic(t string, h TopicHandler) Builder {
  return b.Invoke(h(t, b.fileName))
}

func (b *beebAsm) Do() error {
  return CloseBuilder(b.err, b.w, b.fileName, b.modified)
}

func (b *beebAsm) write(s string) Builder {
  Write(&b.err, &b.w, s)
  return b
}

func (b *beebAsm) append(c, f string, a ...interface{}) Builder {
  s := fmt.Sprintf(f, a...)
  if c != "" {
    s = fmt.Sprintf("%-32s ; %s", s, c)
  }
  return b.write(s)
}

func (b *beebAsm) Comment(s string, a ...interface{}) Builder {
  // This form of comment starts at the beginning of the line
  return b.write(fmt.Sprintf("; "+s, a...))
}

func (b *beebAsm) Header(label, value, comment string) Builder {
  // Value with no label is invalid so ignore
  if label != "" && value != "" {
    return b.append(comment, "%s = %s", label, value)
  }
  return b
}

func (b *beebAsm) Newline() Builder {
  return b.write("")
}

func (b *beebAsm) Separator() Builder {
  return b.write("; " + strings.Repeat("*", 75))
}

func (b *beebAsm) Hex(v interface{}) string {
  if i, ok := util.DecodeInt(v, 0); ok {
    return fmt.Sprintf("&%X", i)
  }
  panic(fmt.Errorf("invalid value %v", v))
}
