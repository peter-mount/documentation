package latex

import (
	"context"
	"golang.org/x/net/html"
)

func head(n *html.Node, ctx context.Context) error {
	return nil
}

func beginDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `% Test LaTeX textbook

\documentclass{textbook}
\usepackage{multirow}
\usepackage{array}

\lang      {english}
\title     {Showcasing The Textbook Class}
\subtitle  {Bring Your Focus Back on Writing}
\author    {Peter Mount}
\license   {CC}{by-nc-sa}{3.0}{The Company}
\isbn      {978-0201529838}
\publisher {NV Publishing Company}
\edition   {1}{2023}
\dedicate  {Myself}{Because I'm cool}
\thank     {Thank you to me for being the best}
\keywords  {academic, template, packages}

\begin{document}`)
}

func endDocument(n *html.Node, ctx context.Context) error {
	return WriteStringLn(ctx, `\end{document}`)
}
