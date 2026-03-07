package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"strconv"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"
	"github.com/FalkZ/md-slides/theming"
	"github.com/FalkZ/md-slides/widgets"
)

var logger *log.Logger

func initLog() {
	f, err := os.Create("slides.log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot create log file:", err)
		os.Exit(1)
	}
	logger = log.New(f, "", log.Ltime|log.Lmicroseconds)
}

func isEmptyCell(cell vaxis.Cell) bool {
	return cell.Character.Grapheme == "" &&
		cell.Style.Foreground == 0 &&
		cell.Style.Background == 0 &&
		cell.Style.Attribute == 0
}

func renderSurface(s vxfw.Surface, win vaxis.Window, cache *imageCache, rootBg vaxis.Color) {
	for i, cell := range s.Buffer {
		if isEmptyCell(cell) {
			continue
		}
		if cell.Style.Background == 0 && rootBg != 0 {
			cell.Style.Background = rootBg
		}
		r := i / int(s.Size.Width)
		c := i % int(s.Size.Width)
		win.SetCell(c, r, cell)
	}

	if iw, ok := s.Widget.(*widgets.ImageWidget); ok {
		cache.draw(iw, win)
	}

	for _, child := range s.Children {
		cw := min(int(child.Surface.Size.Width), int(s.Size.Width)-child.Origin.Col)
		ch := min(int(child.Surface.Size.Height), int(s.Size.Height)-child.Origin.Row)
		if cw <= 0 || ch <= 0 {
			continue
		}
		childWin := win.New(child.Origin.Col, child.Origin.Row, cw, ch)
		renderSurface(child.Surface, childWin, cache, rootBg)
	}
}

func main() {
	debug := false
	dumpPage := -1
	dumpRawPage := -1
	args := []string{}
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--debug":
			debug = true
		case "--schema":
			os.Stdout.Write(theming.ThemeSchema())
			return
		case "--dump", "--dump-raw":
			flag := os.Args[i]
			i++
			if i >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "error: %s requires a page number\n", flag)
				os.Exit(1)
			}
			n, err := strconv.Atoi(os.Args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: invalid page number %q\n", os.Args[i])
				os.Exit(1)
			}
			if flag == "--dump" {
				dumpPage = n
			} else {
				dumpRawPage = n
			}
		default:
			args = append(args, os.Args[i])
		}
	}
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: slides [--debug] [--schema] [--dump <page>] [--dump-raw <page>] <file.md>")
		os.Exit(1)
	}

	initLog()
	widgets.SetLogger(logger)
	logger.Println("starting slides")

	mdPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(mdPath)
	logger.Printf("markdown: %s (base: %s)", mdPath, baseDir)

	raw, err := os.ReadFile(mdPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	currentMode := theming.Dark
	parsed := parseMarkdown(raw, currentMode)
	slides := parsed.slides
	theme := parsed.theme
	if len(slides) == 0 {
		fmt.Fprintln(os.Stderr, "no slides found")
		os.Exit(1)
	}
	logger.Printf("parsed %d slides", len(slides))
	if len(parsed.warnings) > 0 {
		for _, w := range parsed.warnings {
			logger.Printf("theme warning: %s", w)
		}
	}

	if dumpPage >= 0 || dumpRawPage >= 0 {
		dp := dumpPage
		if dp < 0 {
			dp = dumpRawPage
		}
		if dp < 1 || dp > len(slides) {
			fmt.Fprintf(os.Stderr, "error: page %d out of range [1, %d]\n", dp, len(slides))
			os.Exit(1)
		}
		idx := dp - 1
		w, h := 80, 24
		slideW := widgets.NewRootWidget(slides[idx], w, h)
		drawCtx := vxfw.DrawContext{
			Characters: vaxis.Characters,
			Min:        vxfw.Size{Width: uint16(w), Height: uint16(h)},
			Max:        vxfw.Size{Width: uint16(w), Height: uint16(h)},
		}
		surf, err := slideW.Draw(drawCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: draw failed: %v\n", err)
			os.Exit(1)
		}
		if dumpRawPage >= 0 {
			os.Stdout.Write(flattenSurface(surf, w, h))
		} else if x, ok := slideW.(widgets.Xmler); ok {
			os.Stdout.Write(marshalSlideXML(x, idx))
			os.Stdout.Write([]byte("\n"))
		}
		return
	}

	page := 0

	vx, err := vaxis.New(vaxis.Options{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer vx.Close()

	logger.Printf("vaxis initialized — kitty=%v sixel=%v graphics=%v rgb=%v",
		vx.CanKittyGraphics(), vx.CanSixel(), vx.CanDisplayGraphics(), vx.CanRGB())

	cache := newImageCache(vx, baseDir)
	defer cache.clear()
	w, _ := vx.Window().Size()
	contentW := w - 4
	if contentW < 10 {
		contentW = 10
	}
	cache.build(slides)

	drawCtx := vxfw.DrawContext{
		Characters: vaxis.Characters,
	}

	statusStyle := theme.ForMode(currentMode).Status

	render := func() {
		win := vx.Window()
		win.Clear()
		rootBg := theme.ForMode(currentMode).Root
		if rootBg.Background != 0 {
			win.Fill(vaxis.Cell{Style: vaxis.Style{Background: rootBg.Background}})
		}
		w, h := win.Size()
		logger.Printf("render: page=%d window=%dx%d", page, w, h)

		s := slides[page]

		slideW := widgets.NewRootWidget(s, w, h)
		drawCtx.Min = vxfw.Size{Width: uint16(w), Height: uint16(h)}
		drawCtx.Max = vxfw.Size{Width: uint16(w), Height: uint16(h)}

		surf, err := slideW.Draw(drawCtx)
		if err != nil {
			logger.Printf("draw error: %v", err)
			return
		}

		renderSurface(surf, win, cache, rootBg.Background)

		if debug {
			if x, ok := slideW.(widgets.Xmler); ok {
				writeDebugXML(x, page)
			}
		}

		status := fmt.Sprintf("%d / %d", page+1, len(slides))
		if len(parsed.warnings) > 0 {
			status = fmt.Sprintf(" [!%d theme warnings] %s", len(parsed.warnings), status)
		}
		statusCol := w - len(status) - 1
		if statusCol < 0 {
			statusCol = 0
		}
		sub := win.New(statusCol, h-1, len(status)+1, 1)
		st := statusStyle
		if st.Background == 0 && rootBg.Background != 0 {
			st.Background = rootBg.Background
		}
		sub.Print(vaxis.Segment{Text: status, Style: st})

		vx.Render()
	}

	render()

	for ev := range vx.Events() {
		needsRender := false

		switch ev := ev.(type) {
		case vaxis.Key:
			switch {
			case ev.MatchString("Ctrl+c"), ev.MatchString("q"):
				return
			case ev.MatchString("l"), ev.MatchString("Right"), ev.MatchString("n"), ev.MatchString("Space"):
				if page < len(slides)-1 {
					page++
					needsRender = true
				}
			case ev.MatchString("h"), ev.MatchString("Left"), ev.MatchString("p"):
				if page > 0 {
					page--
					needsRender = true
				}
			case ev.MatchString("g"), ev.MatchString("Home"):
				if page != 0 {
					page = 0
					needsRender = true
				}
			case ev.MatchString("G"), ev.MatchString("End"):
				if page != len(slides)-1 {
					page = len(slides) - 1
					needsRender = true
				}
			case ev.MatchString("t"):
				if currentMode == theming.Dark {
					currentMode = theming.Light
				} else {
					currentMode = theming.Dark
				}
				slides = rebuildSlides(parsed.doc, parsed.source, theme, currentMode)
				statusStyle = theme.ForMode(currentMode).Status
				cache.build(slides)
				needsRender = true
				logger.Printf("toggled theme to mode=%d", currentMode)
			}
		case vaxis.ColorThemeUpdate:
			newMode := theming.Dark
			if ev.Mode == vaxis.LightMode {
				newMode = theming.Light
			}
			if newMode != currentMode {
				currentMode = newMode
				slides = rebuildSlides(parsed.doc, parsed.source, theme, currentMode)
				statusStyle = theme.ForMode(currentMode).Status
				cache.build(slides)
				needsRender = true
				logger.Printf("system theme changed to mode=%d", currentMode)
			}
		case vaxis.Resize:
			w, _ := vx.Window().Size()
			contentW := w - 4
			if contentW < 10 {
				contentW = 10
			}
			cache.build(slides)
			needsRender = true
		case vaxis.Redraw:
			render()
		}

		if needsRender {
			render()
		}
	}
}
