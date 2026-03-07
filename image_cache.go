package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"git.sr.ht/~rockorager/vaxis"
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
			resolved := resolvePath(c.baseDir, iw.Path)
			goImg := loadImage(resolved)
			if goImg == nil {
				continue
			}
			vimg, err := c.vx.NewImage(goImg)
			if err != nil {
				logger.Printf("  NewImage error for %s: %v", resolved, err)
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
		logger.Printf("  image open error: %v", err)
		return nil
	}
	defer f.Close()

	goImg, err := png.Decode(f)
	if err != nil {
		f.Seek(0, 0)
		goImg, _, err = image.Decode(f)
		if err != nil {
			logger.Printf("  image decode error: %v", err)
			return nil
		}
	}
	return goImg
}
