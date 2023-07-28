package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
	"strings"
)

func (c *Converter) idiomatic(n *html.Node, ctx context.Context) error {
	if parser.HasClass(n, "fas") {
		// convert the class fa-arrow-left to \faArrowLeft and ignore any content
		for _, class := range parser.GetClass(n) {
			if strings.HasPrefix(class, "fa-") {
				s := strings.Split(class, "-")
				for i := 1; i < len(s); i++ {
					a := s[i]
					a = strings.ToTitle(string(rune(a[0]))) + a[1:]
					s[i] = a
				}
				return util.Writef(ctx, `\%s`, strings.Join(s, ""))
			}
		}
	}

	// Default to the original purpose of <i> - italic
	return util.HandleSimpleCommandSpace(`\textit`, n, ctx)
}
