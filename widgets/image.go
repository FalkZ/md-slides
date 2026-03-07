package widgets

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
)

type ImageWidget struct {
	Path       string
	AltText    string
	CellWidth  int
	CellHeight int
	Style      vaxis.Style
}

func (w *ImageWidget) Xml() XmlNode {
	return XmlNode{
		XMLName: xml.Name{Local: "ImageWidget"},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "path"}, Value: w.Path},
			{Name: xml.Name{Local: "altText"}, Value: w.AltText},
			{Name: xml.Name{Local: "cellWidth"}, Value: fmt.Sprint(w.CellWidth)},
			{Name: xml.Name{Local: "cellHeight"}, Value: fmt.Sprint(w.CellHeight)},
		},
	}
}

func (w *ImageWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	surf := vxfw.NewSurface(uint16(w.CellWidth), uint16(w.CellHeight), w)
	return surf, nil
}

func parseInt(value string, defaultValue int) int {

	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}
	return v
}

func altText(img *ast.Image, source []byte) string {
	var buf = ""
	for c := img.FirstChild(); c != nil; c = c.NextSibling() {
		buf += CollectText(c, source)
	}
	return buf
}

const defaultCellHeight = 20

func NewImageWidget(image *ast.Image, source []byte, theme theming.ModeTheme) *ImageWidget {
	alt := altText(image, source)

	contentW := 40
	parts := strings.SplitN(alt, "|", 2)
	alt = parts[0]
	var cellH = defaultCellHeight
	if len(parts) > 1 {
		cellH = parseInt(parts[1], defaultCellHeight)
	}

	return &ImageWidget{Path: string(image.Destination), AltText: alt, CellWidth: contentW, CellHeight: cellH, Style: theme.Image}
}
