package theming

import (
	"strings"

	"git.sr.ht/~rockorager/vaxis"
)

func ParseClasses(classes string) vaxis.Style {
	var style vaxis.Style
	for _, class := range strings.Fields(classes) {
		switch {
		case class == "font-bold":
			style.Attribute |= vaxis.AttrBold
		case class == "font-italic":
			style.Attribute |= vaxis.AttrItalic
		case class == "font-dim":
			style.Attribute |= vaxis.AttrDim
		case class == "underline", class == "decoration-solid":
			style.UnderlineStyle = vaxis.UnderlineSingle
		case class == "decoration-double":
			style.UnderlineStyle = vaxis.UnderlineDouble
		case class == "decoration-wavy":
			style.UnderlineStyle = vaxis.UnderlineCurly
		case class == "decoration-dotted":
			style.UnderlineStyle = vaxis.UnderlineDotted
		case class == "decoration-dashed":
			style.UnderlineStyle = vaxis.UnderlineDashed
		case class == "no-underline":
			style.UnderlineStyle = vaxis.UnderlineOff
		case strings.HasPrefix(class, "bg-"):
			if color, ok := lookupColor(class[3:]); ok {
				style.Background = color
			}
		case strings.HasPrefix(class, "text-"):
			if color, ok := lookupColor(class[5:]); ok {
				style.Foreground = color
			}
		}
	}
	return style
}

func ValidateClass(class string) bool {
	switch {
	case class == "font-bold", class == "font-italic", class == "font-dim":
		return true
	case class == "underline", class == "no-underline",
		class == "decoration-solid", class == "decoration-double",
		class == "decoration-wavy", class == "decoration-dotted",
		class == "decoration-dashed":
		return true
	case strings.HasPrefix(class, "bg-"):
		_, ok := lookupColor(class[3:])
		return ok
	case strings.HasPrefix(class, "text-"):
		_, ok := lookupColor(class[5:])
		return ok
	default:
		return false
	}
}
