package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/latex/util"
	"github.com/peter-mount/documentation/tools/genlatex/parser"
	"golang.org/x/net/html"
)

const (
	StyleChapter       = `\chapter`
	StyleSection       = `\section`
	StyleSubSection    = `\subsection`
	StyleSubSubSection = `\subsubsection`
	StyleParagraph     = `\paragraph`
)

// heading handles headings h1 h2 etc. based on class type
func (c *Converter) headingStart(n *html.Node, ctx context.Context) error {
	var err error

	hType := n.Data

	// Get the style, note ordering is important here
	// as we want class to override the element name
	var style string
	switch {
	// title is for a non-numbered heading?
	case parser.HasClass(n, "chapter"),
		parser.HasClass(n, "title"):
		style = StyleChapter

	case parser.HasClass(n, "section"):
		style = StyleSection

	case parser.HasClass(n, "subsection"):
		style = StyleSubSection

	case parser.HasClass(n, "subsubsection"):
		style = StyleSubSubSection

	case parser.HasClass(n, "paragraph"):
		style = StyleParagraph

	case hType == "h1":
		style = StyleSection

	case hType == "h2":
		style = StyleSubSection

	case hType == "h3":
		style = StyleSubSubSection

	case hType == "h4", hType == "h5":
		style = StyleParagraph

		// Default numbered heading
	default:
		style = StyleSection
	}
	err = util.WriteString(ctx, style)

	if err == nil {
		err = util.Write(ctx, '{')
	}

	return err
}

func (c *Converter) headingEnd(n *html.Node, ctx context.Context) error {
	return util.Write(ctx, '}', '\n')
}
