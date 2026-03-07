package widgets

import (
	"encoding/xml"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/theming"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

type TableWidget struct {
	HasHeader bool
	Rows      [][]string
	Theme     theming.TableTheme
}

func NewTable(n ast.Node, source []byte, theme theming.ModeTheme) *TableWidget {
	w := &TableWidget{Theme: theme.Table}
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() != east.KindTableHeader && child.Kind() != east.KindTableRow {
			continue
		}
		if !w.HasHeader && child.Kind() == east.KindTableHeader {
			w.HasHeader = true
		}
		var cells []string
		for cell := child.FirstChild(); cell != nil; cell = cell.NextSibling() {
			if cell.Kind() == east.KindTableCell {
				cells = append(cells, CollectText(cell, source))
			}
		}
		w.Rows = append(w.Rows, cells)
	}
	return w
}

func (w *TableWidget) Xml() XmlNode {
	var children []XmlNode
	for i, row := range w.Rows {
		var rowAttrs []xml.Attr
		if i == 0 && w.HasHeader {
			rowAttrs = []xml.Attr{{Name: xml.Name{Local: "header"}, Value: "true"}}
		}
		cells := make([]XmlNode, len(row))
		for j, text := range row {
			cells[j] = XmlNode{
				XMLName: xml.Name{Local: "cell"},
				Text:    text,
			}
		}
		children = append(children, XmlNode{
			XMLName:  xml.Name{Local: "row"},
			Attrs:    rowAttrs,
			Children: cells,
		})
	}
	return XmlNode{XMLName: xml.Name{Local: "TableWidget"}, Children: children}
}

func (w *TableWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	if len(w.Rows) == 0 {
		return vxfw.NewSurface(0, 0, w), nil
	}

	var maxColumns int
	for _, row := range w.Rows {
		if len(row) > maxColumns {
			maxColumns = len(row)
		}
	}

	columnWidths := make([]int, maxColumns)
	for _, row := range w.Rows {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	totalWidth := 0
	for i, width := range columnWidths {
		totalWidth += width
		if i > 0 {
			totalWidth += 3
		}
	}

	totalHeight := len(w.Rows)
	if w.HasHeader {
		totalHeight++
	}
	surface := vxfw.NewSurface(uint16(totalWidth), uint16(totalHeight), w)

	headerStyle := w.Theme.Header
	cellStyle := w.Theme.Cell
	separatorStyle := w.Theme.Separator

	surfaceRow := 0
	for rowIndex, row := range w.Rows {
		style := cellStyle
		if rowIndex == 0 && w.HasHeader {
			style = headerStyle
		}
		column := 0
		for columnIndex := 0; columnIndex < maxColumns; columnIndex++ {
			if columnIndex > 0 {
				for j, character := range []rune{' ', '│', ' '} {
					surface.WriteCell(uint16(column+j), uint16(surfaceRow), vaxis.Cell{
						Character: vaxis.Character{Grapheme: string(character), Width: 1},
						Style:     separatorStyle,
					})
				}
				column += 3
			}
			value := ""
			if columnIndex < len(row) {
				value = row[columnIndex]
			}
			padded := value + strings.Repeat(" ", columnWidths[columnIndex]-len(value))
			for _, character := range padded {
				surface.WriteCell(uint16(column), uint16(surfaceRow), vaxis.Cell{
					Character: vaxis.Character{Grapheme: string(character), Width: 1},
					Style:     style,
				})
				column++
			}
		}
		surfaceRow++

		if rowIndex == 0 && w.HasHeader {
			column = 0
			for columnIndex, width := range columnWidths {
				if columnIndex > 0 {
					for _, character := range []rune{'─', '┼', '─'} {
						surface.WriteCell(uint16(column), uint16(surfaceRow), vaxis.Cell{
							Character: vaxis.Character{Grapheme: string(character), Width: 1},
							Style:     separatorStyle,
						})
						column++
					}
				}
				for j := 0; j < width; j++ {
					surface.WriteCell(uint16(column), uint16(surfaceRow), vaxis.Cell{
						Character: vaxis.Character{Grapheme: "─", Width: 1},
						Style:     separatorStyle,
					})
					column++
				}
			}
			surfaceRow++
		}
	}

	return surface, nil
}
