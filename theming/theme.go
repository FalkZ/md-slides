package theming

import "git.sr.ht/~rockorager/vaxis"

type Mode int

const (
	Dark Mode = iota
	Light
)

type Theme struct {
	Light ModeTheme
	Dark  ModeTheme
}

func (t Theme) ForMode(mode Mode) ModeTheme {
	if mode == Light {
		return t.Light
	}
	return t.Dark
}

type ModeTheme struct {
	Root           vaxis.Style
	Paragraph      vaxis.Style
	Heading1       vaxis.Style
	Heading2       vaxis.Style
	Heading3       vaxis.Style
	Heading4       vaxis.Style
	Heading5       vaxis.Style
	Heading6       vaxis.Style
	Blockquote     vaxis.Style
	Link           vaxis.Style
	CodeInline     vaxis.Style
	CodeBlock      vaxis.Style
	Table          TableTheme
	Image vaxis.Style
	List           vaxis.Style
	Status         vaxis.Style
	Syntax         SyntaxTheme
}

func (m ModeTheme) HeadingStyle(level int) vaxis.Style {
	switch level {
	case 1:
		return m.Heading1
	case 2:
		return m.Heading2
	case 3:
		return m.Heading3
	case 4:
		return m.Heading4
	case 5:
		return m.Heading5
	case 6:
		return m.Heading6
	default:
		return m.Heading1
	}
}

type TableTheme struct {
	Header    vaxis.Style
	Cell      vaxis.Style
	Separator vaxis.Style
}

type SyntaxTheme struct {
	Keyword         vaxis.Style
	String          vaxis.Style
	Comment         vaxis.Style
	Function        vaxis.Style
	FunctionMethod  vaxis.Style
	FunctionBuiltin vaxis.Style
	FunctionName    vaxis.Style
	Type            vaxis.Style
	Number          vaxis.Style
	Operator        vaxis.Style
	Variable        vaxis.Style
	Property        vaxis.Style
	ConstantBuiltin vaxis.Style
	Escape          vaxis.Style
}

func (s SyntaxTheme) StyleFor(capture string) (vaxis.Style, bool) {
	switch capture {
	case "keyword":
		return s.Keyword, true
	case "string":
		return s.String, true
	case "comment":
		return s.Comment, true
	case "function":
		return s.Function, true
	case "function.method":
		return s.FunctionMethod, true
	case "function.builtin":
		return s.FunctionBuiltin, true
	case "function.name":
		return s.FunctionName, true
	case "type":
		return s.Type, true
	case "number":
		return s.Number, true
	case "operator":
		return s.Operator, true
	case "variable":
		return s.Variable, true
	case "property":
		return s.Property, true
	case "constant.builtin":
		return s.ConstantBuiltin, true
	case "escape":
		return s.Escape, true
	default:
		return vaxis.Style{}, false
	}
}
