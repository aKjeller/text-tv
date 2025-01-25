package display

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
)

type cell struct {
	char rune
	bg   color.Color
	fg   color.Color
}

func toColorGrid(g grid, sp svttext.SubPage) grid {
	colors, err := newColorMap(sp.Image, len(g[0]), len(g))
	if err != nil {
		log.Fatalf("failed to create colormap", err)
	}

	for i, row := range g {
		for j, c := range row {
			g[i][j] = colors.getColor(j, i, c.char)
		}
	}

	return g
}

var black = color.RGBA{0, 0, 0, 255}

func (c cell) colorRune() string {
	bg, fg := "", ""

	if c.bg != nil && c.bg != black {
		bgR, bgG, bgB, _ := c.bg.RGBA()
		bg = fmt.Sprintf("\u001b[48;2;%d;%d;%dm", bgR>>8, bgG>>8, bgB>>8)
	}
	if c.fg != nil {
		fgR, fgG, fgB, _ := c.fg.RGBA()
		fg = fmt.Sprintf("\u001b[38;2;%d;%d;%dm", fgR>>8, fgG>>8, fgB>>8)
	}

	return bg + fg + string(c.char) + "\u001b[0m"
}

type colorMap struct {
	image image.Image
	dx    int
	dy    int
}

func newColorMap(base64gif string, width, height int) (*colorMap, error) {
	gifReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64gif))
	img, err := gif.Decode(gifReader)
	if err != nil {
		return nil, fmt.Errorf("error during decoding gif: %w", err)
	}

	return &colorMap{
		image: img,
		dx:    img.Bounds().Dx() / width,
		dy:    img.Bounds().Dy() / height,
	}, nil
}

func (c *colorMap) getColor(x, y int, char rune) cell {
	colors := make(map[color.Color]int)
	for i := x * c.dx; i < x*c.dx+c.dx; i++ {
		for j := y * c.dy; j < y*c.dy+c.dy; j++ {
			colors[c.image.At(i, j)]++
		}
	}

	var bg, fg color.Color
	var bgCount, fgCount int

	for color, count := range colors {
		if count > bgCount {
			fgCount = bgCount
			fg = bg
			bgCount = count
			bg = color
		} else if count > fgCount {
			fgCount = count
			fg = color
		}
	}

	return cell{char: char, bg: bg, fg: fg}
}
