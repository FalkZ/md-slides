package widgets

import (
	"encoding/xml"
	"fmt"

	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/vxlayout"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
	gfmAst "github.com/yuin/goldmark/extension/ast"
)

type gap struct{}

func (gap) Xml() XmlNode {
	return XmlNode{XMLName: xml.Name{Local: "gap"}}
}

func (gap) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	return vxfw.NewSurface(ctx.Max.Width, 1, gap{}), nil
}

type Slide struct {
	Widgets []vxfw.Widget
}

func NewSlide(nodes []ast.Node, source []byte, theme theming.ModeTheme) Slide {
	var widgets []vxfw.Widget
	for _, node := range nodes {
		var widget vxfw.Widget
		switch node.Kind() {
		case ast.KindHeading:
			widget = NewHeading(node.(*ast.Heading), source, theme)
		case ast.KindParagraph:
			if node.ChildCount() == 1 && node.FirstChild().Kind() == ast.KindImage {
				widget = NewImageWidget(node.FirstChild().(*ast.Image), source, theme)
			} else {
				widget = NewParagraph(node, source, theme)
			}
		case ast.KindFencedCodeBlock, ast.KindCodeBlock:
			widget = NewCodeBlock(node, source, theme)
		case ast.KindList:
			widget = NewList(node.(*ast.List), source, theme)
		case ast.KindBlockquote:
			widget = NewBlockquote(node, source, theme)
		default:
			if node.Kind() == gfmAst.KindTable {
				widget = NewTable(node, source, theme)
			}
		}
		if widget != nil {
			widgets = append(widgets, widget)
		}
	}
	return Slide{Widgets: widgets}
}

type rootWidget struct {
	content      *vxlayout.FlexLayout
	windowWidth  int
	windowHeight int
}

func NewRootWidget(s Slide, windowWidth, windowHeight int) vxfw.Widget {
	var items []*vxlayout.FlexItem
	var prev vxfw.Widget
	for _, widget := range s.Widgets {
		if prev != nil {
			_, prevBQ := prev.(*BlockquoteWidget)
			_, curBQ := widget.(*BlockquoteWidget)
			_, prevList := prev.(*ListWidget)
			_, curList := widget.(*ListWidget)
			if !(prevBQ && curBQ) && !(prevList && curList) {
				items = append(items, &vxlayout.FlexItem{Widget: gap{}})
			}
		}
		items = append(items, &vxlayout.FlexItem{Widget: widget})
		prev = widget
	}
	items = append(items, vxlayout.MustSpacer(1))
	return &rootWidget{
		content:      &vxlayout.FlexLayout{Children: items, Direction: vxlayout.FlexVertical},
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
	}
}

func (w *rootWidget) Xml() XmlNode {
	var children []XmlNode
	for _, item := range w.content.Children {
		if x, ok := item.Widget.(Xmler); ok {
			children = append(children, x.Xml())
		}
	}
	return XmlNode{
		XMLName: xml.Name{Local: "rootWidget"},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "windowWidth"}, Value: fmt.Sprint(w.windowWidth)},
			{Name: xml.Name{Local: "windowHeight"}, Value: fmt.Sprint(w.windowHeight)},
		},
		Children: children,
	}
}

func (w *rootWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	childCtx := ctx.WithConstraints(
		vxfw.Size{Width: uint16(w.windowWidth - 4), Height: 0},
		vxfw.Size{Width: uint16(w.windowWidth - 4), Height: uint16(w.windowHeight - 2)},
	)
	childSurf, err := w.content.Draw(childCtx)
	if err != nil {
		return vxfw.Surface{}, err
	}

	surf := vxfw.NewSurface(uint16(w.windowWidth), uint16(w.windowHeight), w)
	surf.AddChild(2, 1, childSurf)
	return surf, nil
}
