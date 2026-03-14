package widgets

import (
	"encoding/xml"
	"fmt"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/richtext"
	"git.sr.ht/~rockorager/vaxis/vxfw/vxlayout"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
)

type ListItemBlock struct {
	Segments []vaxis.Segment
	Sublist  *ListWidget
}

type ListItem struct {
	Blocks []ListItemBlock
}

type ListWidget struct {
	Ordered bool
	Start   int
	Depth   int
	Items   []ListItem
	style   vaxis.Style
}

func newListAtDepth(l *ast.List, source []byte, theme theming.ModeTheme, depth int) *ListWidget {
	w := &ListWidget{
		Ordered: l.IsOrdered(),
		Start:   l.Start,
		Depth:   depth,
		style:   theme.List,
	}
	for item := l.FirstChild(); item != nil; item = item.NextSibling() {
		if item.Kind() != ast.KindListItem {
			continue
		}
		var blocks []ListItemBlock
		for block := item.FirstChild(); block != nil; block = block.NextSibling() {
			if block.Kind() == ast.KindList {
				blocks = append(blocks, ListItemBlock{Sublist: newListAtDepth(block.(*ast.List), source, theme, depth+1)})
			} else if block.Kind() == ast.KindParagraph || block.Kind() == ast.KindTextBlock {
				blocks = append(blocks, ListItemBlock{Segments: collectInlineSegments(block, source)})
			}
		}
		w.Items = append(w.Items, ListItem{Blocks: blocks})
	}
	return w
}

func NewList(l *ast.List, source []byte, theme theming.ModeTheme) *ListWidget {
	return newListAtDepth(l, source, theme, 0)
}

func (w *ListWidget) Xml() XmlNode {
	attrs := []xml.Attr{
		{Name: xml.Name{Local: "ordered"}, Value: fmt.Sprint(w.Ordered)},
		{Name: xml.Name{Local: "depth"}, Value: fmt.Sprint(w.Depth)},
		{Name: xml.Name{Local: "start"}, Value: fmt.Sprint(w.Start)},
	}
	var items []XmlNode
	for _, item := range w.Items {
		var itemChildren []XmlNode
		for _, block := range item.Blocks {
			if block.Sublist != nil {
				itemChildren = append(itemChildren, block.Sublist.Xml())
			} else {
				for _, seg := range block.Segments {
					itemChildren = append(itemChildren, segToXml(seg))
				}
			}
		}
		items = append(items, XmlNode{XMLName: xml.Name{Local: "item"}, Children: itemChildren})
	}
	return XmlNode{XMLName: xml.Name{Local: "ListWidget"}, Attrs: attrs, Children: items}
}

func (w *ListWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	var items []*vxlayout.FlexItem
	indent := strings.Repeat("  ", w.Depth)
	idx := w.Start
	if idx == 0 {
		idx = 1
	}

	for _, item := range w.Items {
		var bullet string
		if w.Ordered {
			bullet = fmt.Sprintf("%s%d. ", indent, idx)
			idx++
		} else {
			bullet = indent + "• "
		}

		first := true
		for _, block := range item.Blocks {
			if block.Sublist != nil {
				items = append(items, &vxlayout.FlexItem{Widget: block.Sublist})
				continue
			}
			segs := make([]vaxis.Segment, len(block.Segments))
			for i, seg := range block.Segments {
				segs[i] = vaxis.Segment{Text: seg.Text, Style: mergeStyle(w.style, seg.Style)}
			}
			if first {
				segs = append([]vaxis.Segment{{Text: bullet, Style: w.style}}, segs...)
				first = false
			} else {
				segs = append([]vaxis.Segment{{Text: indent + "    ", Style: w.style}}, segs...)
			}
			rt := richtext.New(segs)
			items = append(items, &vxlayout.FlexItem{Widget: rt})
		}
		if first {
			rt := richtext.New([]vaxis.Segment{{Text: bullet, Style: w.style}})
			items = append(items, &vxlayout.FlexItem{Widget: rt})
		}
	}

	ctx.Min.Height = 0
	fl := &vxlayout.FlexLayout{Children: items, Direction: vxlayout.FlexVertical}
	return fl.Draw(ctx)
}
