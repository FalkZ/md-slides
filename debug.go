package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/widgets"
)

func marshalSlideXML(root widgets.Xmler, page int) []byte {
	node := root.Xml()
	slide := widgets.XmlNode{
		XMLName:  xml.Name{Local: "slide"},
		Attrs:    []xml.Attr{{Name: xml.Name{Local: "page"}, Value: fmt.Sprint(page)}},
		Children: []widgets.XmlNode{node},
	}
	data, _ := xml.MarshalIndent(slide, "", "  ")
	return append([]byte(xml.Header), data...)
}

func writeDebugXML(root widgets.Xmler, page int) {
	os.WriteFile("slides-debug.xml", marshalSlideXML(root, page), 0644)
}

func flattenSurface(s vxfw.Surface, w, h int) []byte {
	grid := make([][]rune, h)
	for r := range grid {
		grid[r] = make([]rune, w)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	var blit func(s vxfw.Surface, offCol, offRow int)
	blit = func(s vxfw.Surface, offCol, offRow int) {
		sw := int(s.Size.Width)
		for i, cell := range s.Buffer {
			g := cell.Character.Grapheme
			if g == "" {
				continue
			}
			r := offRow + i/sw
			c := offCol + i%sw
			if r >= 0 && r < h && c >= 0 && c < w {
				grid[r][c] = []rune(g)[0]
			}
		}
		for _, child := range s.Children {
			blit(child.Surface, offCol+child.Origin.Col, offRow+child.Origin.Row)
		}
	}
	blit(s, 0, 0)

	lines := make([]string, h)
	for r, row := range grid {
		lines[r] = strings.TrimRight(string(row), " ")
	}

	// trim trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return []byte(strings.Join(lines, "\n") + "\n")
}
