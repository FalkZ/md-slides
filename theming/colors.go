package theming

import "git.sr.ht/~rockorager/vaxis"

func rgb(r, g, b uint8) vaxis.Color {
	return vaxis.RGBColor(r, g, b)
}

var colorPalette = map[string][11]vaxis.Color{
	"slate": {
		rgb(248, 250, 252), rgb(241, 245, 249), rgb(226, 232, 240), rgb(203, 213, 225), rgb(148, 163, 184),
		rgb(100, 116, 139), rgb(71, 85, 105), rgb(51, 65, 85), rgb(30, 41, 59), rgb(15, 23, 42), rgb(2, 6, 23),
	},
	"gray": {
		rgb(249, 250, 251), rgb(243, 244, 246), rgb(229, 231, 235), rgb(209, 213, 219), rgb(156, 163, 175),
		rgb(107, 114, 128), rgb(75, 85, 99), rgb(55, 65, 81), rgb(31, 41, 55), rgb(17, 24, 39), rgb(3, 7, 18),
	},
	"zinc": {
		rgb(250, 250, 250), rgb(244, 244, 245), rgb(228, 228, 231), rgb(212, 212, 216), rgb(161, 161, 170),
		rgb(113, 113, 122), rgb(82, 82, 91), rgb(63, 63, 70), rgb(39, 39, 42), rgb(24, 24, 27), rgb(9, 9, 11),
	},
	"neutral": {
		rgb(250, 250, 250), rgb(245, 245, 245), rgb(229, 229, 229), rgb(212, 212, 212), rgb(163, 163, 163),
		rgb(115, 115, 115), rgb(82, 82, 82), rgb(64, 64, 64), rgb(38, 38, 38), rgb(23, 23, 23), rgb(10, 10, 10),
	},
	"stone": {
		rgb(250, 250, 249), rgb(245, 245, 244), rgb(231, 229, 228), rgb(214, 211, 209), rgb(168, 162, 158),
		rgb(120, 113, 108), rgb(87, 83, 78), rgb(68, 64, 60), rgb(41, 37, 36), rgb(28, 25, 23), rgb(12, 10, 9),
	},
	"red": {
		rgb(254, 242, 242), rgb(254, 226, 226), rgb(254, 202, 202), rgb(252, 165, 165), rgb(248, 113, 113),
		rgb(239, 68, 68), rgb(220, 38, 38), rgb(185, 28, 28), rgb(153, 27, 27), rgb(127, 29, 29), rgb(69, 10, 10),
	},
	"orange": {
		rgb(255, 247, 237), rgb(255, 237, 213), rgb(254, 215, 170), rgb(253, 186, 116), rgb(251, 146, 60),
		rgb(249, 115, 22), rgb(234, 88, 12), rgb(194, 65, 12), rgb(154, 52, 18), rgb(124, 45, 18), rgb(67, 20, 7),
	},
	"amber": {
		rgb(255, 251, 235), rgb(254, 243, 199), rgb(253, 230, 138), rgb(252, 211, 77), rgb(251, 191, 36),
		rgb(245, 158, 11), rgb(217, 119, 6), rgb(180, 83, 9), rgb(146, 64, 14), rgb(120, 53, 15), rgb(69, 26, 3),
	},
	"yellow": {
		rgb(254, 252, 232), rgb(254, 249, 195), rgb(254, 240, 138), rgb(253, 224, 71), rgb(250, 204, 21),
		rgb(234, 179, 8), rgb(202, 138, 4), rgb(161, 98, 7), rgb(133, 77, 14), rgb(113, 63, 18), rgb(66, 32, 6),
	},
	"lime": {
		rgb(247, 254, 231), rgb(236, 252, 203), rgb(217, 249, 157), rgb(190, 242, 100), rgb(163, 230, 53),
		rgb(132, 204, 22), rgb(101, 163, 13), rgb(77, 124, 15), rgb(63, 98, 18), rgb(54, 83, 20), rgb(26, 46, 5),
	},
	"green": {
		rgb(240, 253, 244), rgb(220, 252, 231), rgb(187, 247, 208), rgb(134, 239, 172), rgb(74, 222, 128),
		rgb(34, 197, 94), rgb(22, 163, 74), rgb(21, 128, 61), rgb(22, 101, 52), rgb(20, 83, 45), rgb(5, 46, 22),
	},
	"emerald": {
		rgb(236, 253, 245), rgb(209, 250, 229), rgb(167, 243, 208), rgb(110, 231, 183), rgb(52, 211, 153),
		rgb(16, 185, 129), rgb(5, 150, 105), rgb(4, 120, 87), rgb(6, 95, 70), rgb(6, 78, 59), rgb(2, 44, 34),
	},
	"teal": {
		rgb(240, 253, 250), rgb(204, 251, 241), rgb(153, 246, 228), rgb(94, 234, 212), rgb(45, 212, 191),
		rgb(20, 184, 166), rgb(13, 148, 136), rgb(15, 118, 110), rgb(17, 94, 89), rgb(19, 78, 74), rgb(4, 47, 46),
	},
	"cyan": {
		rgb(236, 254, 255), rgb(207, 250, 254), rgb(165, 243, 252), rgb(103, 232, 249), rgb(34, 211, 238),
		rgb(6, 182, 212), rgb(8, 145, 178), rgb(14, 116, 144), rgb(21, 94, 117), rgb(22, 78, 99), rgb(8, 51, 68),
	},
	"sky": {
		rgb(240, 249, 255), rgb(224, 242, 254), rgb(186, 230, 253), rgb(125, 211, 252), rgb(56, 189, 248),
		rgb(14, 165, 233), rgb(2, 132, 199), rgb(3, 105, 161), rgb(7, 89, 133), rgb(12, 74, 110), rgb(8, 47, 73),
	},
	"blue": {
		rgb(239, 246, 255), rgb(219, 234, 254), rgb(191, 219, 254), rgb(147, 197, 253), rgb(96, 165, 250),
		rgb(59, 130, 246), rgb(37, 99, 235), rgb(29, 78, 216), rgb(30, 64, 175), rgb(30, 58, 138), rgb(23, 37, 84),
	},
	"indigo": {
		rgb(238, 242, 255), rgb(224, 231, 255), rgb(199, 210, 254), rgb(165, 180, 252), rgb(129, 140, 248),
		rgb(99, 102, 241), rgb(79, 70, 229), rgb(67, 56, 202), rgb(55, 48, 163), rgb(49, 46, 129), rgb(30, 27, 75),
	},
	"violet": {
		rgb(245, 243, 255), rgb(237, 233, 254), rgb(221, 214, 254), rgb(196, 181, 253), rgb(167, 139, 250),
		rgb(139, 92, 246), rgb(124, 58, 237), rgb(109, 40, 217), rgb(91, 33, 182), rgb(76, 29, 149), rgb(46, 16, 101),
	},
	"purple": {
		rgb(250, 245, 255), rgb(243, 232, 255), rgb(233, 213, 255), rgb(216, 180, 254), rgb(192, 132, 252),
		rgb(168, 85, 247), rgb(147, 51, 234), rgb(126, 34, 206), rgb(107, 33, 168), rgb(88, 28, 135), rgb(59, 7, 100),
	},
	"fuchsia": {
		rgb(253, 244, 255), rgb(250, 232, 255), rgb(245, 208, 254), rgb(240, 171, 252), rgb(232, 121, 249),
		rgb(217, 70, 239), rgb(192, 38, 211), rgb(162, 28, 175), rgb(134, 25, 143), rgb(112, 26, 117), rgb(74, 4, 78),
	},
	"pink": {
		rgb(253, 242, 248), rgb(252, 231, 243), rgb(251, 207, 232), rgb(249, 168, 212), rgb(244, 114, 182),
		rgb(236, 72, 153), rgb(219, 39, 119), rgb(190, 24, 93), rgb(157, 23, 77), rgb(131, 24, 67), rgb(69, 10, 37),
	},
	"rose": {
		rgb(255, 241, 242), rgb(255, 228, 230), rgb(254, 205, 211), rgb(253, 164, 175), rgb(251, 113, 133),
		rgb(244, 63, 94), rgb(225, 29, 72), rgb(190, 18, 60), rgb(159, 18, 57), rgb(136, 19, 55), rgb(76, 5, 25),
	},
}

var shadeIndex = map[string]int{
	"50": 0, "100": 1, "200": 2, "300": 3, "400": 4,
	"500": 5, "600": 6, "700": 7, "800": 8, "900": 9, "950": 10,
}

var specialColors = map[string]vaxis.Color{
	"white":       rgb(255, 255, 255),
	"black":       rgb(0, 0, 0),
	"transparent": vaxis.Color(0),
}

func lookupColor(name string) (vaxis.Color, bool) {
	if c, ok := parseHexColor(name); ok {
		return c, true
	}
	if c, ok := specialColors[name]; ok {
		return c, true
	}
	parts := splitLast(name)
	if parts == nil {
		return 0, false
	}
	colorName, shade := parts[0], parts[1]
	palette, ok := colorPalette[colorName]
	if !ok {
		return 0, false
	}
	idx, ok := shadeIndex[shade]
	if !ok {
		return 0, false
	}
	return palette[idx], true
}

func parseHexColor(name string) (vaxis.Color, bool) {
	if len(name) < 4 || name[0] != '[' || name[len(name)-1] != ']' {
		return 0, false
	}
	hex := name[1 : len(name)-1]
	if len(hex) != 7 || hex[0] != '#' {
		return 0, false
	}
	r, ok1 := hexByte(hex[1], hex[2])
	g, ok2 := hexByte(hex[3], hex[4])
	b, ok3 := hexByte(hex[5], hex[6])
	if !ok1 || !ok2 || !ok3 {
		return 0, false
	}
	return rgb(r, g, b), true
}

func hexByte(hi, lo byte) (uint8, bool) {
	h, ok1 := hexDigit(hi)
	l, ok2 := hexDigit(lo)
	return h<<4 | l, ok1 && ok2
}

func hexDigit(b byte) (uint8, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}

func splitLast(s string) []string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '-' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}
