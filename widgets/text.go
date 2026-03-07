package widgets

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/richtext"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

var defaultInlineTheme = struct {
	CodeInline vaxis.Style
	Link       vaxis.Style
}{}

type TextWidget struct {
	Segments []vaxis.Segment
	Softwrap bool
}

func NewText(n ast.Node, source []byte) TextWidget {
	return TextWidget{Segments: collectInlineSegments(n, source), Softwrap: true}
}

func NewStyledText(n ast.Node, source []byte, baseStyle vaxis.Style) TextWidget {
	segs := collectInlineSegments(n, source)
	for i := range segs {
		segs[i].Style = mergeStyle(baseStyle, segs[i].Style)
	}
	return TextWidget{Segments: segs, Softwrap: true}
}

func mergeStyle(base, inline vaxis.Style) vaxis.Style {
	result := base
	if inline.Foreground != 0 {
		result.Foreground = inline.Foreground
	}
	if inline.Background != 0 {
		result.Background = inline.Background
	}
	result.Attribute |= inline.Attribute
	if inline.UnderlineStyle != 0 {
		result.UnderlineStyle = inline.UnderlineStyle
	}
	if inline.Hyperlink != "" {
		result.Hyperlink = inline.Hyperlink
	}
	return result
}

func (w TextWidget) Xml() XmlNode {
	return XmlNode{
		XMLName: xml.Name{Local: "TextWidget"},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "softwrap"}, Value: fmt.Sprint(w.Softwrap)},
		},
		Children: w.SegmentsXml(),
	}
}

func (w TextWidget) SegmentsXml() []XmlNode {
	nodes := make([]XmlNode, len(w.Segments))
	for i, seg := range w.Segments {
		nodes[i] = segToXml(seg)
	}
	return nodes
}

func (w TextWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	rt := richtext.New(w.Segments)
	rt.Softwrap = w.Softwrap
	return rt.Draw(ctx)
}

func CollectText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if child.Kind() == ast.KindText {
			buf.Write(child.(*ast.Text).Value(source))
		} else if child.Kind() == ast.KindString {
			buf.Write(child.(*ast.String).Value)
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}

func blockText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		buf.Write(seg.Value(source))
	}
	return buf.String()
}

func SetInlineTheme(theme theming.ModeTheme) {
	defaultInlineTheme.CodeInline = theme.CodeInline
	defaultInlineTheme.Link = theme.Link
}

func collectInlineSegments(n ast.Node, source []byte) []vaxis.Segment {
	var segs []vaxis.Segment
	var styleStack []vaxis.Style

	currentStyle := func() vaxis.Style {
		if len(styleStack) > 0 {
			return styleStack[len(styleStack)-1]
		}
		return vaxis.Style{}
	}

	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		switch child.Kind() {
		case ast.KindEmphasis:
			em := child.(*ast.Emphasis)
			if entering {
				s := currentStyle()
				if em.Level == 1 {
					s.Attribute |= vaxis.AttrItalic
				} else {
					s.Attribute |= vaxis.AttrBold
				}
				styleStack = append(styleStack, s)
			} else if len(styleStack) > 0 {
				styleStack = styleStack[:len(styleStack)-1]
			}

		case ast.KindCodeSpan:
			if entering {
				styleStack = append(styleStack, defaultInlineTheme.CodeInline)
			} else if len(styleStack) > 0 {
				styleStack = styleStack[:len(styleStack)-1]
			}

		case ast.KindLink:
			if entering {
				link := child.(*ast.Link)
				s := currentStyle()
				s.Foreground = defaultInlineTheme.Link.Foreground
				s.Attribute |= defaultInlineTheme.Link.Attribute
				s.Hyperlink = string(link.Destination)
				styleStack = append(styleStack, s)
			} else if len(styleStack) > 0 {
				styleStack = styleStack[:len(styleStack)-1]
			}

		case east.KindStrikethrough:
			if entering {
				s := currentStyle()
				s.Attribute |= vaxis.AttrStrikethrough
				styleStack = append(styleStack, s)
			} else if len(styleStack) > 0 {
				styleStack = styleStack[:len(styleStack)-1]
			}

		case ast.KindText:
			if entering {
				t := child.(*ast.Text)
				text := strings.TrimRight(string(t.Value(source)), "\n")
				if t.SoftLineBreak() {
					text += " "
				}
				segs = append(segs, vaxis.Segment{
					Text:  text,
					Style: currentStyle(),
				})
			}

		case ast.KindString:
			if entering {
				s := child.(*ast.String)
				segs = append(segs, vaxis.Segment{
					Text:  string(s.Value),
					Style: currentStyle(),
				})
			}

		case ast.KindAutoLink:
			if entering {
				al := child.(*ast.AutoLink)
				s := currentStyle()
				s.Foreground = defaultInlineTheme.Link.Foreground
				s.Attribute |= defaultInlineTheme.Link.Attribute
				s.Hyperlink = string(al.URL(source))
				segs = append(segs, vaxis.Segment{
					Text:  string(al.Label(source)),
					Style: s,
				})
			}
			return ast.WalkSkipChildren, nil

		case ast.KindImage:
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return mergeSegments(segs)
}

func mergeSegments(segs []vaxis.Segment) []vaxis.Segment {
	if len(segs) <= 1 {
		return segs
	}
	merged := []vaxis.Segment{segs[0]}
	for _, seg := range segs[1:] {
		last := &merged[len(merged)-1]
		if last.Style == seg.Style {
			last.Text += seg.Text
		} else {
			merged = append(merged, seg)
		}
	}
	return merged
}
