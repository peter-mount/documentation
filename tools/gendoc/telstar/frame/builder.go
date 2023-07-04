package frame

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/documentation/tools/gendoc/telstar/edittf"
	"strconv"
	"strings"
)

// Builder allows for the creation of content
type Builder struct {
	content []byte // Frame content
	cx, cy  int    // Cursor position
	wrap    bool   // Wrap at end of line or stop?
}

const (
	frameWidth     = 40
	frameHeight    = 25
	frameWidthMax  = frameWidth - 1
	frameHeightMax = frameHeight - 1
)

func Of(frame []string) *Builder {
	var buf bytes.Buffer
	for _, l := range edittf.PadLines(frame) {
		buf.Write([]byte(l))
	}
	return &Builder{content: buf.Bytes()}
}

// Builder returns a Builder based on this frame
func (f *Frame) Builder() *Builder {
	switch f.Content.Type {
	case "edit.tf":
		return Of(edittf.Decode(f.Content.Data))
	case "markup":
		// TODO decode markup first
		return Of(strings.Split(f.Content.Data, "\r\n"))
	case "raw":
		// TODO decode raw first
		return Of(strings.Split(f.Content.Data, "\r\n"))
	default:
		return Of(strings.Split(f.Content.Data, "\r\n"))
	}
}

func (b *Builder) Wrap(v string) *Builder {
	b.wrap = true
	return b
}

func (b *Builder) NoWrap(v string) *Builder {
	b.wrap = false
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minMax(a, b, c int) int {
	return min(a, max(b, c))
}

func offset(x, y int) int {
	return (frameWidth * y) + x
}

func (b *Builder) isValid(x, y int) bool {
	return b.isXValid(x) && b.isYValid(y)
}

func (b *Builder) isXValid(x int) bool {
	return x >= 0 && x < frameWidth
}

func (b *Builder) isYValid(y int) bool {
	return y >= 0 && y < frameHeight
}

// SetPos sets the current text cursor in the frame
func (b *Builder) SetPos(x, y int) *Builder {
	if b.isXValid(x) {
		b.cx = x
	}
	if b.isYValid(y) {
		b.cy = y
	}
	return b
}

func (b *Builder) GetChar(x, y int) byte {
	if b.isValid(x, y) {
		return b.content[offset(x, y)]
	}
	return 0
}

// Set sets the character at a specific position on the frame. This does not use the text cursor
func (b *Builder) Set(x, y int, c byte) *Builder {
	if b.isValid(x, y) {
		b.content[offset(x, y)] = c
	}
	return b
}

// LF alias of Down
func (b *Builder) LF() *Builder {
	return b.Down()
}

// Down moves the text cursor down 1 character.
// If wrap is enabled then when the bottom is reached it will wrap to the top of the frame.
func (b *Builder) Down() *Builder {
	if b.cy < 23 {
		b.cy = b.cy + 1
	} else if b.wrap {
		b.cy = 0
	}
	return b
}

// Up moves the text cursor up 1 character.
// If wrap is enabled then when the top is reached it will wrap to the bottom of the frame.
func (b *Builder) Up() *Builder {
	if b.cy > 0 {
		b.cy = b.cy - 1
	} else if b.wrap {
		b.cy = 23
	}
	return b
}

// Right moves the text cursor right 1 character.
// If wrap is enabled then when the right is reached it will wrap to the left of the frame and the next line down.
func (b *Builder) Right() *Builder {
	if b.cx < 39 {
		b.cx = b.cx + 1
	} else if b.wrap {
		return b.CR().Down()
	}
	return b
}

// Left moves the text cursor left 1 character.
// If wrap is enabled then when the left is reached it will wrap to the right of the frame and the next line up
func (b *Builder) Left() *Builder {
	if b.cx > 0 {
		b.cx = b.cx - 1
	} else if b.wrap {
		b.cx = 39
		b.Up()
	}
	return b
}

// CR moves the text cursor to the start of the current line.
func (b *Builder) CR() *Builder {
	b.cx = 0
	return b
}

// Newline moves the text cursor to the start of the next line. This is the same as CR().LF()
func (b *Builder) Newline() *Builder {
	return b.CR().LF()
}

// WriteBytes writes the byte slice to the frame, moving the cursor as applicable
func (b *Builder) WriteBytes(v ...byte) *Builder {
	for _, c := range v {
		switch c {
		case 8:
			b.Left()
		case 9:
			b.Right()
		case 10:
			b.Down()
		case 11:
			b.Up()
		case 13:
			b.CR()
		default:
			if c > 31 {
				b.Set(b.cx, b.cy, c).Right()
			}
		}
	}
	return b
}

// Write writes the string to the frame, moving the cursor as applicable
func (b *Builder) Write(v string) *Builder {
	return b.WriteBytes([]byte(v)...)
}

// Writef writes the formatted string to the frame, moving the cursor as applicable
func (b *Builder) Writef(v string, a ...interface{}) *Builder {
	return b.Write(fmt.Sprintf(v, a...))
}

// WriteInt writes the integer to the frame, moving the cursor as applicable
func (b *Builder) WriteInt(v int) *Builder {
	return b.Write(strconv.Itoa(v))
}
