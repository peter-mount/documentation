package custom

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	util2 "github.com/peter-mount/documentation/tools/gensite/latex/util"
	"golang.org/x/net/html"
)

func OpcodeTable6502(n *html.Node, ctx context.Context) error {
	err := util.WriteString(ctx, "\n\\asmOpcodes{")
	if err == nil {
		err = opcodeTableBody(n, ctx)
	}
	if err == nil {
		err = util.WriteString(ctx, "}")
	}
	return err
}

func opcodeTableBody(n *html.Node, ctx context.Context) error {
	tRow := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "thead":
				if err := util.WriteString(ctx, "\\asmOpcodeHeader{}"); err != nil {
					return err
				}
				// Embed any tableAlign marginNote
				notes := util.BuffersFromContext(ctx).GetBytes("marginNote")
				if notes != nil {
					if err := util.Write(ctx, notes...); err != nil {
						return err
					}
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

	case parser.HasClass(td, "opcodeHex"):
		return util.HandleSimpleCommand(`\asmOpcodeCell`, td, ctx)

	case parser.HasClasses(td, "opcodeBytes", "opcodeCycles"):
		var left, right string
		for c := td.FirstChild; c != nil; c = c.NextSibling {
			switch c.Type {
			case html.TextNode:
				left = c.Data
			case html.ElementNode:
				if c.Data == "sup" {
					right = parser.GetText(c)
				}
			}
		}
		return util.Writef(ctx, `\asmOpcodeWithNote{%s}{%s}`, util2.EscapeText(left), util2.EscapeText(right))

	case parser.HasClass(td, "processorSupported"):
		return util.WriteString(ctx, `\asmOpcodeSupported`)

	case parser.HasClass(td, "processorUnsupported"):
		return util.WriteString(ctx, `\asmOpcodeUnsupported`)

	default:
		return nil
	}
}
