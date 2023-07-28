package custom

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func TruthTable(n *html.Node, ctx context.Context) error {
	fmt := `\truthTable{%s}{%s}{%s}{%s}{%s}`
	if parser.HasClass(n, "marginNote") {
		fmt = `\marginnote{` + fmt + "}"
	}
	return util.Writef(ctx, fmt,
		parser.GetTextByClass(n, "truthTableCaption", ""),
		parser.GetTextByClass(n, "truthTable00", "~"),
		parser.GetTextByClass(n, "truthTable01", "~"),
		parser.GetTextByClass(n, "truthTable10", "~"),
		parser.GetTextByClass(n, "truthTable11", "~"),
	)
}
