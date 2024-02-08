package util

import (
	"bytes"
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

type Buffers struct {
	m map[string]*bytes.Buffer
}

func (b *Buffers) GetBuffer(n string) *bytes.Buffer {
	if v, exists := b.m[n]; exists {
		return v
	}
	return &bytes.Buffer{}
}

func (b *Buffers) GetOrCreateBuffer(n string) *bytes.Buffer {
	buffer := b.m[n]
	if buffer == nil {
		buffer = &bytes.Buffer{}
		b.m[n] = buffer
	}
	return buffer
}

func (b *Buffers) NewBuffer(n string) *bytes.Buffer {
	buffer := &bytes.Buffer{}
	b.m[n] = buffer
	return buffer
}

func (b *Buffers) End(n string) {
	delete(b.m, n)
}

func (b *Buffers) UseBuffer(name string, h parser.Handler, n *html.Node, ctx context.Context) error {
	buffer := b.GetOrCreateBuffer(name)
	return h.Do(n, WithContext(buffer, ctx))
}

func (b *Buffers) GetBytes(n string) []byte {
	if buffer, exists := b.m[n]; exists {
		r := buffer.Bytes()
		b.End(n)
		return r
	}
	return nil
}

func (b *Buffers) GetString(n string) string {
	a := b.GetBytes(n)
	if a == nil {
		return ""
	}
	return string(a)
}

func NewBuffers(ctx context.Context) context.Context {
	b := &Buffers{m: make(map[string]*bytes.Buffer)}
	return context.WithValue(ctx, "util.buffers", b)
}

func BuffersFromContext(ctx context.Context) *Buffers {
	b := ctx.Value("util.buffers")
	if b != nil {
		return b.(*Buffers)
	}
	return nil
}
