package widgets

import (
	"encoding/xml"
	"fmt"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
)

type XmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",attr"`
	Children []XmlNode  `xml:",any"`
	Text     string     `xml:",chardata"`
}

type Xmler interface {
	Xml() XmlNode
}

func segToXml(seg vaxis.Segment) XmlNode {
	attrs := []xml.Attr{{Name: xml.Name{Local: "text"}, Value: seg.Text}}
	if seg.Style.Foreground != 0 {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "fg"}, Value: fmt.Sprint(seg.Style.Foreground)})
	}
	if seg.Style.Background != 0 {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "bg"}, Value: fmt.Sprint(seg.Style.Background)})
	}
	if seg.Style.Attribute != 0 {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "attr"}, Value: attrString(seg.Style.Attribute)})
	}
	if seg.Style.Hyperlink != "" {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "hyperlink"}, Value: seg.Style.Hyperlink})
	}
	return XmlNode{XMLName: xml.Name{Local: "seg"}, Attrs: attrs}
}

func attrString(a vaxis.AttributeMask) string {
	var parts []string
	if a&vaxis.AttrBold != 0 {
		parts = append(parts, "bold")
	}
	if a&vaxis.AttrDim != 0 {
		parts = append(parts, "dim")
	}
	if a&vaxis.AttrItalic != 0 {
		parts = append(parts, "italic")
	}
	if a&vaxis.AttrBlink != 0 {
		parts = append(parts, "blink")
	}
	if a&vaxis.AttrReverse != 0 {
		parts = append(parts, "reverse")
	}
	if a&vaxis.AttrInvisible != 0 {
		parts = append(parts, "invisible")
	}
	if a&vaxis.AttrStrikethrough != 0 {
		parts = append(parts, "strikethrough")
	}
	return strings.Join(parts, "|")
}
