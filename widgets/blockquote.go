package widgets

import (
	"encoding/xml"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/richtext"
	"git.sr.ht/~rockorager/vaxis/vxfw/vxlayout"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
)

type BlockquoteWidget struct {
	Paragraphs []TextWidget
}

func NewBlockquote(n ast.Node, source []byte, theme theming.ModeTheme) *BlockquoteWidget {
	style := theme.Blockquote
	var paragraphs []TextWidget
	for block := n.FirstChild(); block != nil; block = block.NextSibling() {
		if block.Kind() == ast.KindParagraph {
			tw := NewStyledText(block, source, style)
			tw.Segments = append([]vaxis.Segment{{Text: "▍ ", Style: style}}, tw.Segments...)
			paragraphs = append(paragraphs, tw)
		}
	}
	return &BlockquoteWidget{Paragraphs: paragraphs}
}

func (w *BlockquoteWidget) Xml() XmlNode {
	var children []XmlNode
	for _, paragraph := range w.Paragraphs {
		children = append(children, XmlNode{
			XMLName:  xml.Name{Local: "paragraph"},
			Children: paragraph.SegmentsXml(),
		})
	}
	return XmlNode{
		XMLName:  xml.Name{Local: "BlockquoteWidget"},
		Children: children,
	}
}

func (w *BlockquoteWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	var items []*vxlayout.FlexItem
	for _, paragraph := range w.Paragraphs {
		rt := richtext.New(paragraph.Segments)
		rt.Softwrap = paragraph.Softwrap
		items = append(items, &vxlayout.FlexItem{Widget: rt})
	}
	ctx.Min.Height = 0
	fl := &vxlayout.FlexLayout{Children: items, Direction: vxlayout.FlexVertical}
	return fl.Draw(ctx)
}
