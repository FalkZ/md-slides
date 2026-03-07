package widgets

import (
	"encoding/xml"
	"fmt"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	vxtext "git.sr.ht/~rockorager/vaxis/vxfw/text"
	"github.com/FalkZ/md-slides/theming"
)

type HRWidget struct {
	Width int
	Style vaxis.Style
}

func NewHR(width int, theme theming.ModeTheme) *HRWidget {
	return &HRWidget{Width: width, Style: theme.HorizontalRule}
}

func (w *HRWidget) Xml() XmlNode {
	return XmlNode{
		XMLName: xml.Name{Local: "HRWidget"},
		Attrs:   []xml.Attr{{Name: xml.Name{Local: "width"}, Value: fmt.Sprint(w.Width)}},
	}
}

func (w *HRWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	t := vxtext.New(strings.Repeat("─", w.Width))
	t.Style = w.Style
	return t.Draw(ctx)
}
