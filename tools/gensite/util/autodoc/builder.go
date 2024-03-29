package autodoc

import (
	"context"
	"github.com/peter-mount/documentation/tools/gensite"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Builder interface {
	// Comment generates a comment with nothing else on the line
	Comment(string, ...interface{}) Builder
	// Header sets a constant value to a label
	Header(string, string, string) Builder
	// Function declares a function
	Function(string, string, string) Builder
	// Newline writes a blank line
	Newline() Builder
	// Separator writes a Comment line with a row of characters to mark some division in the output
	Separator() Builder
	// InvokeTopic invokes a TopicHandler, e.g. contains additional details about the file being generated
	InvokeTopic(string, TopicHandler) Builder
	// Invoke invokes a Handler
	Invoke(Handler) Builder
	// Using calls a Provider to add a sub Builder to receive calls.
	// Used to generate multiple versions of the same data. Most Builder implementations should panic if this is called.
	Using(p Provider) Builder
	// Do runs the Builder to produce the output file(s)
	Do(ctx context.Context) error
	// Hex convert a value to Hexadecimal in their specific format
	Hex(interface{}) string
}

type Provider func(string, string, time.Time, context.Context) Builder

type Handler func(context.Context, Builder) error

type TopicHandler func(string, string) Handler

// InitBuilder creates any required directories & _index.html files for a built file
// dir directory to write to
// file FileName of file (without suffix)
// asm Assembler - this will form the directory underneath dir where the actual file will live.
// suffix File suffix to add to file, will have "." inserted between them
// title, linkTitle, desc frontmatter for the _index.html for the assembler page
func InitBuilder(dir, file string, modified time.Time, asm, suffix, title, desc string, ctx context.Context) (string, io.WriteCloser, error) {
	if err := GenerateReferenceIndexFile(ctx, path.Join(dir, asm, "_index.html"), modified, title, desc); err != nil {
		return "", nil, err
	}

	writeNow := true
	fsName := file + "." + suffix
	buildFileName := path.Join(dir, asm, fsName)

	if fi, err := os.Stat(buildFileName); err != nil {
		// If error is not-existing then fall through with writeNow is true
		if !os.IsNotExist(err) {
			return "", nil, err
		}

		writeNow = true
	} else {
		writeNow = modified.After(fi.ModTime())
	}

	if err := os.MkdirAll(path.Dir(buildFileName), 0755); err != nil {
		return "", nil, err
	}

	// Add indices files later
	task.GetQueue(ctx).
		AddPriorityTask(tools.PriorityAutodoc, GenerateReferenceIndices(buildFileName, modified).
			WithContext(ctx, ResourceManagerKey)).
		AddPriorityTask(tools.PriorityAutodoc, GenerateFileIndexPage(dir, file, asm, buildFileName, modified).
			WithContext(ctx, ResourceManagerKey))

	if writeNow {
		f, err := os.Create(buildFileName)
		if err != nil {
			return "", nil, err
		}

		return buildFileName, f, nil
	}

	return buildFileName, nil, nil
}

func CloseBuilder(err error, w io.WriteCloser, fileName string, modified time.Time, ctx context.Context) error {
	if w != nil {
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
	ctx      context.Context
	src      []Builder
}

func For(dir, file string, modified time.Time, ctx context.Context) Builder {
	return &unionBuilder{dir: dir, file: file, modified: modified, ctx: ctx}
}

func (u *unionBuilder) Using(p Provider) Builder {
	u.src = append(u.src, p(u.dir, u.file, u.modified, u.ctx))
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

func (u *unionBuilder) Do(ctx context.Context) error {
	for _, b := range u.src {
		err := b.Do(ctx)
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

func (u *unionBuilder) Function(s string, s2 string, s3 string) Builder {
	for _, b := range u.src {
		b.Function(s, s2, s3)
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

func (u *unionBuilder) Hex(interface{}) string {
	panic("not supported")
}
