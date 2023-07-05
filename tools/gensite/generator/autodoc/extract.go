package autodoc

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite/generator"
	"github.com/peter-mount/documentation/tools/gensite/hugo"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/walk"
)

func (s *Autodoc) extract(ctx context.Context) error {
	book := generator.GetBook(ctx)

	// Only run once per Book ID
	if s.extracted.Contains(book.ID) {
		return nil
	}

	// Prevent us running again for this Book
	s.extracted.Add(book.ID)

	log.Printf("Scanning %s for autodocs", book.ID)

	return walk.NewPathWalker().
		IsFile().
		PathNotContain("/reference/").
		PathHasSuffix(".html").
		Then(hugo.FrontMatterActionOf().
			OtherExists("api", s.extractApi).
			OtherExists("memorymap", s.extractMemoryMap).
			Walk(ctx)).
		Walk(book.ContentPath())
}

func (s *Autodoc) extractMemoryMap(ctx context.Context, fm *hugo.FrontMatter) error {

	headers := s.getHeaders(ctx)
	if d, exists := fm.Other["description"]; exists {
		_ = headers.Add(&Header{Comment: util.DecodeString(d, "")})
		log.Println(d)
	}

	return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
		return util.IfMap(e, func(m map[interface{}]interface{}) error {
			var val string
			if e, exists := m["address"]; exists {
				val = "0x" + util.DecodeString(e, "") // Valid as address is always in hex
			} else if e, exists := m["value"]; exists {
				val = util.DecodeString(e, "")
			}
			return headers.Add(&Header{
				Label:   util.DecodeString(m["name"], ""),
				Value:   val,
				Comment: util.DecodeString(m["desc"], ""),
			})
		})
	})
}

func (s *Autodoc) extractApi(ctx context.Context, fm *hugo.FrontMatter) error {

	api := s.GetApi(ctx)
	// TODO add comment here

	return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
		return util.IfMap(e, func(m map[interface{}]interface{}) error {
			v := &ApiEntry{
				Name:     util.DecodeString(m["name"], ""),
				Addr:     "0x" + util.DecodeString(m["addr"], ""),
				Indirect: util.DecodeString(m["indirect"], ""),
				Title:    util.DecodeString(m["title"], ""),
				params:   e,
			}

			if i, ok := util.DecodeInt(v.Addr, 0); !ok {
				return fmt.Errorf("failed to parse addr %q for %s", v.Addr, v.Name)
			} else {
				v.call = i
			}

			if err := util.IfMapEntry(m, "entry", v.Entry.decode); err != nil {
				return err
			}

			if err := util.IfMapEntry(m, "exit", v.Exit.decode); err != nil {
				return err
			}

			api.Add(v)

			return nil
		})
	})
}
