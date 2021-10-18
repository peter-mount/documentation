package asm

import (
  "context"
  "fmt"
  "github.com/peter-mount/documentation/tools/util"
  "github.com/peter-mount/documentation/tools/util/autodoc"
  "io"
  "strings"
  "time"
)

type zAsm struct {
  fileName string
  modified time.Time
  ctx      context.Context
  w        io.WriteCloser
  err      error
}

func (b *zAsm) valid() bool {
  return b.err == nil && b.w != nil
}

func ZAsm(dir, file string, modified time.Time, ctx context.Context) autodoc.Builder {
  fileName, w, err := autodoc.InitBuilder(dir, file, modified, "zasm", "z80", "ZAsm", "Files for the ZAsm assembler", ctx)

  return &zAsm{
    fileName: fileName,
    modified: modified,
    w:        w,
    err:      err,
    ctx:      ctx,
  }
}

func (b *zAsm) Using(autodoc.Provider) autodoc.Builder {
  panic("not implemented")
}

func (b *zAsm) Invoke(handler autodoc.Handler) autodoc.Builder {
  if b.err == nil {
    b.err = handler(b.ctx, b)
  }
  return b
}

func (b *zAsm) InvokeTopic(t string, h autodoc.TopicHandler) autodoc.Builder {
  return b.Invoke(h(t, b.fileName))
}

func (b *zAsm) Do() error {
  return autodoc.CloseBuilder(b.err, b.w, b.fileName, b.modified, b.ctx)
}

func (b *zAsm) write(s string) autodoc.Builder {
  autodoc.Write(&b.err, &b.w, s)
  return b
}

func (b *zAsm) append(c, f string, a ...interface{}) autodoc.Builder {
  s := fmt.Sprintf(f, a...)
  if c != "" {
    s = fmt.Sprintf("%-32s ; %s", s, c)
  }
  return b.write(s)
}

func (b *zAsm) Comment(s string, a ...interface{}) autodoc.Builder {
  // This form of comment starts at the beginning of the line
  return b.write(fmt.Sprintf("; "+s, a...))
}

func (b *zAsm) Header(label, value, comment string) autodoc.Builder {
  // Value with no label is invalid so ignore
  if label != "" && value != "" {
    return b.append(comment, "%-16s equ %s", label, value)
  }
  return b
}

func (b *zAsm) Newline() autodoc.Builder {
  return b.write("")
}

func (b *zAsm) Separator() autodoc.Builder {
  return b.write("; " + strings.Repeat("*", 75))
}

func (b *zAsm) Hex(v interface{}) string {
  if i, ok := util.DecodeInt(v, 0); ok {
    return fmt.Sprintf("&%X", i)
  }
  panic(fmt.Errorf("invalid value %v", v))
}
