package svg

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"math"
)

type DiskDefinition struct {
	Tracks int        `yaml:"tracks"`
	Colour string     `yaml:"colour"`
	Def    []TrackDef `yaml:"trackDefs"`
	radius int
	tStep  int
}

func (d DiskDefinition) getTrackDef(t int) TrackDef {
	td := d.Def[0]
	for _, e := range d.Def {
		if e.Start <= t {
			td = e
		}
	}
	return td
}

func (d DiskDefinition) drawTrack(t int, canvas *svg.SVG) {
	d.getTrackDef(t).
		drawTrack(d.trackRadius(t), canvas)
}

func (d DiskDefinition) trackRadius(t int) int {
	return d.radius - (d.tStep * (t + 2))
}

type TrackDef struct {
	Start   int    `yaml:"start"`
	Sectors int    `yaml:"sectors"`
	Colour  string `yaml:"colour"`
}

func (td TrackDef) drawTrack(r int, canvas *svg.SVG) {
	canvas.Circle(0, 0, r, "fill='"+td.Colour+"'")
}

func (td TrackDef) drawSector(r1, r2 int, canvas *svg.SVG) {
	da := math.Pi * (2.0 / float64(td.Sectors))
	for s := 0; s < td.Sectors; s++ {
		sx, sy := math.Sincos(float64(s) * da)
		canvas.Line(int(float64(r1)*sx), -int(float64(r1)*sy), int(float64(r2)*sx), -int(float64(r2)*sy))
	}
}

func (d DiskDefinition) Write(s *File, canvas *svg.SVG) error {
	d.radius = (s.MinSize() - 10) / 2
	d.tStep = (d.radius - 30) / (d.Tracks)

	// r=total radius in svg, rs=step between tracks, range of tracks 50..250
	w := d.radius + 5

	canvas.Gtransform(fmt.Sprintf("translate(%d %d)", w, w))
	canvas.Group("stroke='black'", "stroke-width='1px'")

	// Disk media
	track := func(r int, c string) {
		canvas.Circle(0, 0, r, "fill='"+c+"'")
	}

	// Disk media
	track(d.radius, d.Colour)

	// Tracks
	for t := 1; t <= d.Tracks; t++ {
		d.drawTrack(t, canvas)
	}

	// Next track if it existed, defines boundary
	track(d.trackRadius(d.Tracks+1), d.Colour)

	// sectors
	t2 := d.Tracks + 1
	for i := len(d.Def) - 1; i >= 0; i-- {
		td2 := d.Def[i]
		td2.drawSector(d.trackRadius(td2.Start), d.trackRadius(t2), canvas)
		t2 = td2.Start
	}

	// Spindle
	track(20, "#fff")

	canvas.Gend()
	canvas.Gend()

	return nil
}
