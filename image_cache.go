package main

import (
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~rockorager/vaxis"
	"golang.org/x/image/draw"
	"github.com/FalkZ/md-slides/widgets"
)

type cachedImage struct {
	img        vaxis.Image
	cellWidth  int
	cellHeight int
}

type imageCache struct {
	vx      *vaxis.Vaxis
	baseDir string
	cache   map[string]cachedImage
}

func newImageCache(vx *vaxis.Vaxis, baseDir string) *imageCache {
	return &imageCache{vx: vx, baseDir: baseDir, cache: map[string]cachedImage{}}
}

func cacheKey(path string, cellHeight int) string {
	return fmt.Sprintf("%s@%d", path, cellHeight)
}

func upscaleImage(img image.Image, cellWidth, cellHeight int) image.Image {
	minWidth := cellWidth * 20
	minHeight := cellHeight * 40
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	scaleX := float64(minWidth) / float64(width)
	scaleY := float64(minHeight) / float64(height)
	scale := scaleX
	if scaleY > scale {
		scale = scaleY
	}
	if scale <= 1.0 {
		return img
	}
	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}

func (c *imageCache) build(slides []widgets.Slide) {
	c.clear()
	seen := map[string]bool{}
	for _, s := range slides {
		for _, w := range s.Widgets {
			iw, ok := w.(*widgets.ImageWidget)
			if !ok {
				continue
			}
			key := cacheKey(iw.Path, iw.CellHeight)
			if seen[key] {
				continue
			}
			seen[key] = true
			var resolved string
			if strings.HasPrefix(iw.Path, "http://") || strings.HasPrefix(iw.Path, "https://") {
				var err error
				resolved, err = resolveUrl(iw.Path)
				if err != nil {
					slog.Debug("image fetch error", "url", iw.Path, "error", err)
					continue
				}
			} else {
				resolved = resolvePath(c.baseDir, iw.Path)
			}
			goImg := loadImage(resolved)
			if goImg == nil {
				continue
			}
			goImg = upscaleImage(goImg, iw.CellWidth, iw.CellHeight)
			vimg, err := c.vx.NewImage(goImg)
			if err != nil {
				slog.Debug("NewImage error", "path", resolved, "error", err)
				continue
			}
			vimg.Resize(iw.CellWidth, iw.CellHeight)
			c.cache[key] = cachedImage{img: vimg, cellWidth: iw.CellWidth, cellHeight: iw.CellHeight}
		}
	}
}

func (c *imageCache) draw(imageWidget *widgets.ImageWidget, win vaxis.Window) {
	key := cacheKey(imageWidget.Path, imageWidget.CellHeight)
	ci, ok := c.cache[key]
	if !ok {
		win.Print(vaxis.Segment{
			Text:  imageWidget.AltText,
			Style: imageWidget.Style,
		})
		return
	}
	ci.img.Draw(win)
}

func (c *imageCache) clear() {
	for k, ci := range c.cache {
		ci.img.Destroy()
		delete(c.cache, k)
	}
}

func resolvePath(base, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(base, p)
}

func loadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		slog.Debug("image open error", "error", err)
		return nil
	}
	defer f.Close()

	goImg, err := png.Decode(f)
	if err != nil {
		f.Seek(0, 0)
		goImg, _, err = image.Decode(f)
		if err != nil {
			slog.Debug("image decode error", "error", err)
			return nil
		}
	}
	return goImg
}
