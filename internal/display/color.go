package display

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"strings"
)

type colorMap struct {
	image image.Image
	dx    int
	dy    int
}

var background = color.RGBA{0, 0, 0, 255}

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

func (c *colorMap) getColor(x, y int) int {
	colors := make(map[color.Color]int)
	for i := x * c.dx; i < x*c.dx+c.dx; i++ {
		for j := y * c.dy; j < y*c.dy+c.dy; j++ {
			colors[c.image.At(i, j)]++
		}
	}

	var color color.Color
	var most int
	for key, value := range colors {
		if key != background {
			if value > most {
				most = value
				color = key
			}
		}
	}

	return toAnsi(color)
}

// TODO use correct colors
func toAnsi(color color.Color) int {
	if color == yellow {
		return 33
	} else if color == blue {
		return 34
	} else if color == green {
		return 32
	} else if color == bg_blue {
		return 42
	}
	return 0
}

var yellow = color.RGBA{255, 255, 0, 255}
var blue = color.RGBA{0, 255, 255, 255}
var green = color.RGBA{0, 255, 0, 255}
var bg_blue = color.RGBA{0, 0, 255, 255}
