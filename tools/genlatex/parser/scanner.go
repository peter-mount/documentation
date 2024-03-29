package parser

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type Scanner struct {
	handlers       map[string]Handler // Map of handlers this scanner handles
	defaultHandler Handler            // Handler to call if an entry in handlers is not found
	text           Handler            // Handler for handling html.TextNode's
	comment        Handler            // Handler for handling html.CommentNode's
	docType        Handler            // Handler for handling html.DoctypeNode's
}

const (
	scannerKey = "html.parser.scanner"
)

func (s *Scanner) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scannerKey, s)
}

func ScannerFromContext(ctx context.Context) *Scanner {
	v := ctx.Value(scannerKey)
	if v == nil {
		return nil
	}
	return v.(*Scanner)
}

// Handle implements Handler, so we can use a Scanner as a Handler
func (s *Scanner) Handle(n *html.Node, ctx context.Context) error {
	// Set context entry if not existing or a different Scanner
	if ps := ScannerFromContext(ctx); ps == nil || ps != s {
		ctx = s.WithContext(ctx)
	}

	switch n.Type {
	case html.DocumentNode:
		return HandleChildren(s.Handle, n, ctx)

	case html.ElementNode:
		// Do not handle non-printable elements
		if HasClass(n, "book-index") {
			if err := s.indexEntry(n, ctx); err != nil {
				return nil
			}
		}
		if !HasClasses(n, "d-print-none", "lead") {
			h, exists := s.handlers[n.Data]
			if !exists {
				//fmt.Printf("default %q %v\n", n.Data, exists)
				h = s.defaultHandler
			}
			return h.Do(n, ctx)
		}

	case html.CommentNode:
		return s.comment.Do(n, ctx)

	case html.DoctypeNode:
		return s.docType.Do(n, ctx)

	case html.TextNode:
		return s.text.Do(n, ctx)
	}

	return nil
}

func (s *Scanner) indexEntry(n *html.Node, ctx context.Context) error {
	a := GetAttr(n, "data-book-index")
	if a == "" {
		return nil
	}
	w := ctx.Value("writer").(io.Writer)
	_, err := w.Write([]byte(fmt.Sprintf(`\index{%s}`, a)))
	return err
}
