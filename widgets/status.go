package widgets

import (
	"encoding/xml"
	"fmt"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
)

type StatusWidget struct {
	Text        string
	WindowWidth int
	Style       vaxis.Style
}

func (w *StatusWidget) Xml() XmlNode {
	return XmlNode{
		XMLName: xml.Name{Local: "StatusWidget"},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "text"}, Value: w.Text},
			{Name: xml.Name{Local: "windowWidth"}, Value: fmt.Sprint(w.WindowWidth)},
		},
	}
}

func (w *StatusWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	style := w.Style
	width := uint16(w.WindowWidth)
	surf := vxfw.NewSurface(width, 1, w)
	col := int(width) - len(w.Text) - 1
	if col < 0 {
		col = 0
	}
	for _, ch := range w.Text {
		surf.WriteCell(uint16(col), 0, vaxis.Cell{
			Character: vaxis.Character{Grapheme: string(ch), Width: 1},
			Style:     style,
		})
		col++
	}
	return surf, nil
}
