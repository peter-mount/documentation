package custom

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

// ProcessorFlags handles the processorFlags table which is not a true table
func ProcessorFlags(n *html.Node, ctx context.Context) error {

	err := util.WriteString(ctx, "\n\n\\processorFlags{%\n")
	if err == nil {
		err = processorFlagsBody(n, ctx)
	}
	if err == nil {
		err = util.WriteString(ctx, "}\n")
	}
	return err
}

func processorFlagsBody(n *html.Node, ctx context.Context) error {
	tRow := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "thead", "tbody":
				return processorFlagsBody(c, ctx)

			case "tr":
				switch tRow {
				case 0:
					if err := processorFlagsTopRow(c, ctx); err != nil {
						return err
					}
				default:
					// TODO implement subsequent rows
				}
			}
			tRow++
		}
	}
	return nil
}

func processorFlagsTopRow(n *html.Node, ctx context.Context) error {

	rNum := 0
	for tr := parser.FindElement(n, "tr"); tr != nil; tr = tr.NextSibling {

		if err := util.WriteString(ctx, `\processorFlagRow{`); err != nil {
			return err
		}

		for c := tr.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				switch c.Data {
				case "th":
					if rNum == 0 {
						if err := util.Writef(ctx, `\processorFlagHeader{%s}`, parser.GetText(c)); err != nil {
							return err
						}
					} else {
						if err := util.Writef(ctx, `\processorFlagLabel{%s}`, parser.GetText(c)); err != nil {
							return err
						}
					}
				case "td":
					switch rNum {
					case 0:
						ctr := parser.FindElement(c, "tr")
						if ctr != nil {
							for fc := ctr.FirstChild; fc != nil; fc = fc.NextSibling {
								if err := util.Writef(ctx, `\processorFlag{%s}`, parser.GetText(fc)); err != nil {
									return err
								}
							}
						}
					default:
						if err := util.Writef(ctx, `\processorFlagDef{%s}`, parser.GetText(c)); err != nil {
							return err
						}
					}
				}
			}
		}

		if err := util.WriteString(ctx, "}\n"); err != nil {
			return err
		}

		rNum++
	}
	return nil
}
