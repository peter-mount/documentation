package svg

import (
	svg "github.com/ajstarks/svgo"
	"github.com/peter-mount/documentation/tools/gendoc/util"
	"io"
)

type File struct {
	Rect       `yaml:",inline"`
	FileName   string          `yaml:"file"`
	FloppyDisk *DiskDefinition `yaml:"floppyDisk"`
}

type Rect struct {
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

func (r *Rect) MinSize() int {
	if r.Width < r.Height {
		return r.Width
	}
	return r.Height
}

type Handler interface {
	Write(*File, *svg.SVG) error
}

func (s *File) failFast(canvas *svg.SVG, handlers ...Handler) error {
	for _, h := range handlers {
		if h != nil {
			err := h.Write(s, canvas)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *File) FileHandler() util.FileHandler {
	return func(w io.Writer) error {
		canvas := svg.New(w)
		canvas.Startview(s.Width, s.Height, 0, 0, s.Width, s.Height)
		defer canvas.End()

		return s.failFast(canvas,
			s.FloppyDisk)
	}
}
