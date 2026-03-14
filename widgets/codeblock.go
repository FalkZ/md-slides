package widgets

import (
	"encoding/xml"
	"log/slog"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"git.sr.ht/~rockorager/vaxis/vxfw/richtext"
	vxtext "git.sr.ht/~rockorager/vaxis/vxfw/text"
	"github.com/FalkZ/md-slides/theming"
	"github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
	"github.com/yuin/goldmark/ast"
)

type CodeBlockWidget struct {
	Code   string
	Lang   string
	Syntax theming.SyntaxTheme
	Style  vaxis.Style
}

func NewCodeBlock(n ast.Node, source []byte, theme theming.ModeTheme) *CodeBlockWidget {
	lang := ""
	if fcb, ok := n.(*ast.FencedCodeBlock); ok {
		lang = string(fcb.Language(source))
	}
	return &CodeBlockWidget{Code: blockText(n, source), Lang: lang, Syntax: theme.Syntax, Style: theme.CodeBlock}
}

func (w *CodeBlockWidget) Xml() XmlNode {
	return XmlNode{
		XMLName: xml.Name{Local: "CodeBlockWidget"},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "code"}, Value: w.Code},
			{Name: xml.Name{Local: "lang"}, Value: w.Lang},
		},
	}
}

func (w *CodeBlockWidget) Draw(ctx vxfw.DrawContext) (vxfw.Surface, error) {
	code := strings.TrimRight(w.Code, "\n")
	var surf vxfw.Surface
	var err error
	if segs := w.highlightCode(code, w.Lang); segs != nil {
		rt := richtext.New(segs)
		surf, err = rt.Draw(ctx)
	} else {
		t := vxtext.New(code)
		t.Style = w.Style
		surf, err = t.Draw(ctx)
	}
	if err != nil {
		return surf, err
	}
	if w.Style.Background != 0 {
		fillBg(&surf, w.Style.Background)
	}
	return surf, nil
}

func fillBg(surf *vxfw.Surface, bg vaxis.Color) {
	for i, cell := range surf.Buffer {
		if cell.Character.Grapheme == "" && cell.Style.Background == 0 {
			surf.Buffer[i] = vaxis.Cell{
				Character: vaxis.Character{Grapheme: " ", Width: 1},
				Style:     vaxis.Style{Background: bg},
			}
		}
	}
}

var highlighterCache = map[string]*gotreesitter.Highlighter{}

func getHighlighter(lang string) *gotreesitter.Highlighter {
	if h, ok := highlighterCache[lang]; ok {
		return h
	}
	var entry *grammars.LangEntry
	for _, e := range grammars.AllLanguages() {
		if e.Name == lang {
			entry = &e
			break
		}
	}
	if entry == nil {
		highlighterCache[lang] = nil
		return nil
	}
	langObj := entry.Language()
	var opts []gotreesitter.HighlighterOption
	if entry.TokenSourceFactory != nil {
		factory := entry.TokenSourceFactory
		l := langObj
		opts = append(opts, gotreesitter.WithTokenSourceFactory(func(src []byte) gotreesitter.TokenSource {
			return factory(src, l)
		}))
	}
	h, err := gotreesitter.NewHighlighter(langObj, entry.HighlightQuery, opts...)
	if err != nil {
		slog.Debug("highlighter error", "lang", lang, "error", err)
		highlighterCache[lang] = nil
		return nil
	}
	highlighterCache[lang] = h
	return h
}

func (w *CodeBlockWidget) highlightCode(code, lang string) []vaxis.Segment {
	h := getHighlighter(lang)
	if h == nil {
		return nil
	}
	src := []byte(code)
	ranges := h.Highlight(src)
	if len(ranges) == 0 {
		return nil
	}
	var segs []vaxis.Segment
	pos := uint32(0)
	for _, r := range ranges {
		if r.StartByte > pos {
			segs = append(segs, vaxis.Segment{
				Text:  string(src[pos:r.StartByte]),
				Style: w.Style,
			})
		}
		style := w.Style
		if s, ok := w.Syntax.StyleFor(r.Capture); ok {
			s.Background = w.Style.Background
			style = s
		}
		segs = append(segs, vaxis.Segment{
			Text:  string(src[r.StartByte:r.EndByte]),
			Style: style,
		})
		pos = r.EndByte
	}
	if pos < uint32(len(src)) {
		segs = append(segs, vaxis.Segment{
			Text:  string(src[pos:]),
			Style: w.Style,
		})
	}
	return segs
}
