package theming

import (
	_ "embed"
	"os"

	"git.sr.ht/~rockorager/vaxis"
	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var mergeBaseYAML []byte

var fallbackThemeYAML []byte

//go:embed theme_schema.json
var themeSchemaJSON []byte

func SetFallbackThemeYAML(data []byte) {
	fallbackThemeYAML = data
}

func ThemeSchema() []byte {
	return themeSchemaJSON
}

type rawTheme struct {
	Base string       `yaml:"base"`
	Light   rawModeTheme `yaml:"light"`
	Dark    rawModeTheme `yaml:"dark"`
}

type rawModeTheme struct {
	Root       string       `yaml:"root"`
	Paragraph  string       `yaml:"p"`
	Heading1   string       `yaml:"h1"`
	Heading2   string       `yaml:"h2"`
	Heading3   string       `yaml:"h3"`
	Heading4   string       `yaml:"h4"`
	Heading5   string       `yaml:"h5"`
	Heading6   string       `yaml:"h6"`
	Blockquote string       `yaml:"blockquote"`
	Link       string       `yaml:"link"`
	Image      string       `yaml:"img"`
	List       string       `yaml:"list"`
	Status string `yaml:"status"`
	Code       rawCodeTheme `yaml:"code"`
	Table      rawTableTheme `yaml:"table"`
}

type rawCodeTheme struct {
	Block  string          `yaml:"block"`
	Inline string          `yaml:"inline"`
	Syntax rawSyntaxTheme  `yaml:"syntax"`
}

type rawSyntaxTheme struct {
	Keyword         string `yaml:"keyword"`
	String          string `yaml:"string"`
	Comment         string `yaml:"comment"`
	Function        string `yaml:"function"`
	FunctionMethod  string `yaml:"function_method"`
	FunctionBuiltin string `yaml:"function_builtin"`
	FunctionName    string `yaml:"function_name"`
	Type            string `yaml:"type"`
	Number          string `yaml:"number"`
	Operator        string `yaml:"operator"`
	Variable        string `yaml:"variable"`
	Property        string `yaml:"property"`
	ConstantBuiltin string `yaml:"constant_builtin"`
	Escape          string `yaml:"escape"`
}

type rawTableTheme struct {
	Header    string `yaml:"header"`
	Cell      string `yaml:"cell"`
	Separator string `yaml:"separator"`
}

func DefaultTheme() Theme {
	return parseRawTheme(mergeBaseYAML)
}

func FallbackTheme() Theme {
	return parseRawTheme(fallbackThemeYAML)
}

func ParseTheme(frontmatterYAML []byte) Theme {
	return ParseThemeWithResolver(frontmatterYAML, nil)
}

func ParseThemeWithResolver(frontmatterYAML []byte, resolver func(string) (string, error)) Theme {
	if len(frontmatterYAML) == 0 {
		return FallbackTheme()
	}
	var wrapper struct {
		Theme *rawTheme `yaml:"theme"`
	}
	if err := yaml.Unmarshal(frontmatterYAML, &wrapper); err != nil {
		return FallbackTheme()
	}
	if wrapper.Theme == nil {
		return FallbackTheme()
	}

	base := DefaultTheme()
	if wrapper.Theme.Base != "" && resolver != nil {
		path, err := resolver(wrapper.Theme.Base)
		if err == nil {
			data, err := os.ReadFile(path)
			if err == nil {
				baseRaw := parseRawTheme(data)
				base = MergeThemes(base, baseRaw)
			}
		}
	}

	override := resolveTheme(*wrapper.Theme)
	return MergeThemes(base, override)
}

func parseRawTheme(data []byte) Theme {
	var raw rawTheme
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Theme{}
	}
	return resolveTheme(raw)
}

func resolveTheme(raw rawTheme) Theme {
	return Theme{
		Light: resolveModeTheme(raw.Light),
		Dark:  resolveModeTheme(raw.Dark),
	}
}

func resolveModeTheme(raw rawModeTheme) ModeTheme {
	return ModeTheme{
		Root:           ParseClasses(raw.Root),
		Paragraph:      ParseClasses(raw.Paragraph),
		Heading1:       ParseClasses(raw.Heading1),
		Heading2:       ParseClasses(raw.Heading2),
		Heading3:       ParseClasses(raw.Heading3),
		Heading4:       ParseClasses(raw.Heading4),
		Heading5:       ParseClasses(raw.Heading5),
		Heading6:       ParseClasses(raw.Heading6),
		Blockquote:     ParseClasses(raw.Blockquote),
		Link:           ParseClasses(raw.Link),
		Image:          ParseClasses(raw.Image),
		List:           ParseClasses(raw.List),
		Status: ParseClasses(raw.Status),
		CodeInline:     ParseClasses(raw.Code.Inline),
		CodeBlock:      ParseClasses(raw.Code.Block),
		Table: TableTheme{
			Header:    ParseClasses(raw.Table.Header),
			Cell:      ParseClasses(raw.Table.Cell),
			Separator: ParseClasses(raw.Table.Separator),
		},
		Syntax: SyntaxTheme{
			Keyword:         ParseClasses(raw.Code.Syntax.Keyword),
			String:          ParseClasses(raw.Code.Syntax.String),
			Comment:         ParseClasses(raw.Code.Syntax.Comment),
			Function:        ParseClasses(raw.Code.Syntax.Function),
			FunctionMethod:  ParseClasses(raw.Code.Syntax.FunctionMethod),
			FunctionBuiltin: ParseClasses(raw.Code.Syntax.FunctionBuiltin),
			FunctionName:    ParseClasses(raw.Code.Syntax.FunctionName),
			Type:            ParseClasses(raw.Code.Syntax.Type),
			Number:          ParseClasses(raw.Code.Syntax.Number),
			Operator:        ParseClasses(raw.Code.Syntax.Operator),
			Variable:        ParseClasses(raw.Code.Syntax.Variable),
			Property:        ParseClasses(raw.Code.Syntax.Property),
			ConstantBuiltin: ParseClasses(raw.Code.Syntax.ConstantBuiltin),
			Escape:          ParseClasses(raw.Code.Syntax.Escape),
		},
	}
}

func MergeThemes(base, override Theme) Theme {
	return Theme{
		Light: mergeModeTheme(base.Light, override.Light),
		Dark:  mergeModeTheme(base.Dark, override.Dark),
	}
}

func mergeModeTheme(base, override ModeTheme) ModeTheme {
	return ModeTheme{
		Root:           mergeStyle(base.Root, override.Root),
		Paragraph:      mergeStyle(base.Paragraph, override.Paragraph),
		Heading1:       mergeStyle(base.Heading1, override.Heading1),
		Heading2:       mergeStyle(base.Heading2, override.Heading2),
		Heading3:       mergeStyle(base.Heading3, override.Heading3),
		Heading4:       mergeStyle(base.Heading4, override.Heading4),
		Heading5:       mergeStyle(base.Heading5, override.Heading5),
		Heading6:       mergeStyle(base.Heading6, override.Heading6),
		Blockquote:     mergeStyle(base.Blockquote, override.Blockquote),
		Link:           mergeStyle(base.Link, override.Link),
		CodeInline:     mergeStyle(base.CodeInline, override.CodeInline),
		CodeBlock:      mergeStyle(base.CodeBlock, override.CodeBlock),
		Image:  mergeStyle(base.Image, override.Image),
		List:           mergeStyle(base.List, override.List),
		Status:         mergeStyle(base.Status, override.Status),
		Table: TableTheme{
			Header:    mergeStyle(base.Table.Header, override.Table.Header),
			Cell:      mergeStyle(base.Table.Cell, override.Table.Cell),
			Separator: mergeStyle(base.Table.Separator, override.Table.Separator),
		},
		Syntax: mergeSyntaxTheme(base.Syntax, override.Syntax),
	}
}

func mergeSyntaxTheme(base, override SyntaxTheme) SyntaxTheme {
	return SyntaxTheme{
		Keyword:         mergeStyle(base.Keyword, override.Keyword),
		String:          mergeStyle(base.String, override.String),
		Comment:         mergeStyle(base.Comment, override.Comment),
		Function:        mergeStyle(base.Function, override.Function),
		FunctionMethod:  mergeStyle(base.FunctionMethod, override.FunctionMethod),
		FunctionBuiltin: mergeStyle(base.FunctionBuiltin, override.FunctionBuiltin),
		FunctionName:    mergeStyle(base.FunctionName, override.FunctionName),
		Type:            mergeStyle(base.Type, override.Type),
		Number:          mergeStyle(base.Number, override.Number),
		Operator:        mergeStyle(base.Operator, override.Operator),
		Variable:        mergeStyle(base.Variable, override.Variable),
		Property:        mergeStyle(base.Property, override.Property),
		ConstantBuiltin: mergeStyle(base.ConstantBuiltin, override.ConstantBuiltin),
		Escape:          mergeStyle(base.Escape, override.Escape),
	}
}

func mergeStyle(base, override vaxis.Style) vaxis.Style {
	result := base
	if override.Foreground != 0 {
		result.Foreground = override.Foreground
	}
	if override.Background != 0 {
		result.Background = override.Background
	}
	if override.Attribute != 0 {
		result.Attribute = override.Attribute
	}
	if override.UnderlineStyle != 0 {
		result.UnderlineStyle = override.UnderlineStyle
	}
	return result
}
