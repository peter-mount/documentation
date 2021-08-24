package bbc

import (
	"github.com/peter-mount/documentation/tools/util"
	"sort"
)

type Osword struct {
	call   int                         // Call 0..255
	params map[interface{}]interface{} // Params
}

func (b *BBC) extractOsword(osword interface{}) {
	if a, ok := osword.([]interface{}); ok {
		for _, e := range a {
			if m, ok := e.(map[interface{}]interface{}); ok {
				if v, ok := util.DecodeInt(m["int"], 0); ok {
					o := &Osword{
						call:   v,
						params: m,
					}
					b.osword = append(b.osword, o)
				}
			}
		}
	}
}

func (b *BBC) writeOswordIndex() error {
	sort.SliceStable(b.osword, func(i, j int) bool {
		return b.osword[i].call < b.osword[j].call
	})

	r := Output{Nometa: true}
	for _, o := range b.osword {
		r.Osword = append(r.Osword, o.params)
	}

	return util.WriteReferenceYaml(
		*b.baseDir,
		"osword",
		"OSWord calls",
		"OSWord &FFF1 calls",
		10,
		r)
}
