package widgets

import (
	"encoding/xml"

	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
)

type ParagraphWidget struct {
	Text TextWidget
}

func NewParagraph(n ast.Node, source []byte, theme theming.ModeTheme) *ParagraphWidget {
	return &ParagraphWidget{Text: NewStyledText(n, source, theme.Paragraph)}
}

func (w *ParagraphWidget) Xml() XmlNode {
	return XmlNode{XMLName: xml.Name{Local: "ParagraphWidget"}, Children: w.Text.SegmentsXml()}
}

func (w *ParagraphWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	return w.Text.Draw(ctx)
}
