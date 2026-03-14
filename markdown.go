package main

import (
	"strings"

	"github.com/FalkZ/md-slides/theming"
	"github.com/FalkZ/md-slides/widgets"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

type parsedMarkdown struct {
	slides   []widgets.Slide
	theme    theming.Theme
	warnings []theming.ValidationWarning
	doc      ast.Node
	source   []byte
}

func parseMarkdown(raw []byte, mode theming.Mode, baseDir string) parsedMarkdown {
	frontmatter, body := theming.ExtractFrontmatter(raw)
	resolver := func(path string) (string, error) {
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			return resolveUrl(path)
		}
		return resolvePath(baseDir, path), nil
	}
	theme := theming.ParseThemeWithResolver(frontmatter, resolver)
	warnings := theming.Validate(frontmatter)
	doc := md.Parser().Parse(text.NewReader(body))
	modeTheme := theme.ForMode(mode)
	widgets.SetInlineTheme(modeTheme)
	slides := splitSlides(doc, body, modeTheme)
	return parsedMarkdown{
		slides:   slides,
		theme:    theme,
		warnings: warnings,
		doc:      doc,
		source:   body,
	}
}

func rebuildSlides(doc ast.Node, source []byte, theme theming.Theme, mode theming.Mode) []widgets.Slide {
	modeTheme := theme.ForMode(mode)
	widgets.SetInlineTheme(modeTheme)
	return splitSlides(doc, source, modeTheme)
}

func splitSlides(doc ast.Node, source []byte, theme theming.ModeTheme) []widgets.Slide {
	var slides []widgets.Slide
	var nodes []ast.Node
	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() == ast.KindThematicBreak {
			if len(nodes) > 0 {
				slides = append(slides, widgets.NewSlide(nodes, source, theme))
				nodes = nil
			}
			continue
		}
		nodes = append(nodes, child)
	}
	if len(nodes) > 0 {
		slides = append(slides, widgets.NewSlide(nodes, source, theme))
	}
	return slides
}
