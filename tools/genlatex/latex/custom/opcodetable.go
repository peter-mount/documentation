package custom

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

func OpcodeTable6502(n *html.Node, ctx context.Context) error {
	err := util.WriteString(ctx, "\n\n\\asmOpcodes{%\n")
	if err == nil {
		err = opcodeTableBody(n, ctx)
	}
	if err == nil {
		err = util.WriteString(ctx, "}\n")
	}
	return err
}

func opcodeTableBody(n *html.Node, ctx context.Context) error {
	tRow := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "thead":
				if err := util.WriteString(ctx, "\\asmOpcodeHeader{}%\n"); err != nil {
					return err
				}

			case "tbody":
				return opcodeTableBody(c, ctx)

			case "tr":
				if err := opcodeTableRow(c, ctx); err != nil {
					return err
				}
			}
			tRow++
		}
	}
	return nil
}

func opcodeTableRow(tr *html.Node, ctx context.Context) error {
	if err := util.WriteString(ctx, `\asmOpcodeRow{`); err != nil {
		return err
	}

	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "td":
				if err := opcodeTableCell(c, ctx); err != nil {
					return nil
				}
			}
		}
	}

	return util.WriteString(ctx, "}\n")
}

func opcodeTableCell(td *html.Node, ctx context.Context) error {
	switch {
	case parser.HasClass(td, "opcodeDef"):
		return util.HandleSimpleCommand(`\asmOpcodeSyntax`, td, ctx)

	case parser.HasClasses(td, "opcodeHex", "opcodeBytes", "opcodeCycles"):
		return util.HandleSimpleCommand(`\asmOpcodeCell`, td, ctx)

	case parser.HasClass(td, "processorSupported"):
		return util.WriteString(ctx, `\asmOpcodeSupported`)

	case parser.HasClass(td, "processorUnsupported"):
		return util.WriteString(ctx, `\asmOpcodeUnsupported`)

	default:
		return nil
	}
}
