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

type beebAsm struct {
  fileName string
  modified time.Time
  ctx      context.Context
  w        io.WriteCloser
  err      error
}

func (b *beebAsm) valid() bool {
  return b.err == nil && b.w != nil
}

func BeebAsm(dir, file string, modified time.Time, ctx context.Context) autodoc.Builder {
  fileName, w, err := autodoc.InitBuilder(dir, file, modified, "beebasm", "asm", "BeebASM", "Files for the BeebASM assembler", ctx)

  return &beebAsm{
    fileName: fileName,
    modified: modified,
    w:        w,
    err:      err,
    ctx:      ctx,
  }
}

func (b *beebAsm) Using(autodoc.Provider) autodoc.Builder {
  panic("not implemented")
}

func (b *beebAsm) Invoke(handler autodoc.Handler) autodoc.Builder {
  if b.err == nil {
    b.err = handler(b.ctx, b)
  }
  return b
}

func (b *beebAsm) InvokeTopic(t string, h autodoc.TopicHandler) autodoc.Builder {
  return b.Invoke(h(t, b.fileName))
}

func (b *beebAsm) Do(_ context.Context) error {
  return autodoc.CloseBuilder(b.err, b.w, b.fileName, b.modified, b.ctx)
}

func (b *beebAsm) write(s string) autodoc.Builder {
  autodoc.Write(&b.err, &b.w, s)
  return b
}

func (b *beebAsm) append(c, f string, a ...interface{}) autodoc.Builder {
  s := fmt.Sprintf(f, a...)
  if c != "" {
    s = fmt.Sprintf("%-32s ; %s", s, c)
  }
  return b.write(s)
}

func (b *beebAsm) Comment(s string, a ...interface{}) autodoc.Builder {
  // This form of comment starts at the beginning of the line
  return b.write(fmt.Sprintf("; "+s, a...))
}

func (b *beebAsm) Header(label, value, comment string) autodoc.Builder {
  // Value with no label is invalid so ignore
  if label != "" && value != "" {
    return b.append(comment, "%s = %s", label, value)
  }
  return b
}

func (b *beebAsm) Function(label, value, comment string) autodoc.Builder {
  // Value with no label is invalid so ignore
  if label != "" && value != "" {
    return b.append(comment, "%s = %s", label, value)
  }
  return b
}

func (b *beebAsm) Newline() autodoc.Builder {
  return b.write("")
}

func (b *beebAsm) Separator() autodoc.Builder {
  return b.write("; " + strings.Repeat("*", 75))
}

func (b *beebAsm) Hex(v interface{}) string {
  if i, ok := util.DecodeInt(v, 0); ok {
    return fmt.Sprintf("&%X", i)
  }
  panic(fmt.Errorf("invalid value %v", v))
}
