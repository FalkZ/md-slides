package widgets

import (
	"encoding/xml"

	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
)

type HeadingWidget struct {
	Text TextWidget
}

func NewHeading(h *ast.Heading, source []byte, theme theming.ModeTheme) *HeadingWidget {
	style := theme.HeadingStyle(h.Level)
	return &HeadingWidget{Text: NewStyledText(h, source, style)}
}

func (w *HeadingWidget) Xml() XmlNode {
	return XmlNode{XMLName: xml.Name{Local: "HeadingWidget"}, Children: w.Text.SegmentsXml()}
}

func (w *HeadingWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	return w.Text.Draw(ctx)
}
