package svttext

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"strings"
	"time"
)

const (
	PageWidth  = 41
	PageHeight = 25
)

var black = color.RGBA{0, 0, 0, 255}

type Page struct {
	PageNumber string    `json:"pageNumber"`
	PrevPage   string    `json:"prevPage"`
	NextPage   string    `json:"nextPage"`
	SubPages   []SubPage `json:"subPages"`
	FetchedAt  time.Time `json:"fetchedAt"`
}

type SubPage struct {
	Number string
	Text   string
	Grid   Grid
	Image  image.Image
}

type Grid []Row

type Row []Cell

type Cell struct {
	Char rune
	Bg   color.Color
	Fg   color.Color
}

func (r Row) ColorString() string {
	var sb strings.Builder
	for _, c := range r {
		sb.WriteString(c.ColorRune())
	}
	return sb.String()
}

func (c Cell) ColorRune() string {
	bg, fg := "", ""

	if c.Bg != nil && c.Bg != black {
		bgR, bgG, bgB, _ := c.Bg.RGBA()
		bg = fmt.Sprintf("\u001b[48;2;%d;%d;%dm", bgR>>8, bgG>>8, bgB>>8)
	}
	if c.Fg != nil {
		fgR, fgG, fgB, _ := c.Fg.RGBA()
		fg = fmt.Sprintf("\u001b[38;2;%d;%d;%dm", fgR>>8, fgG>>8, fgB>>8)
	}

	return bg + fg + string(c.Char) + "\u001b[0m"
}

func newPage(data data) (Page, error) {

	var subPages []SubPage
	for _, page := range data.SubPages {
		sp, err := newSubPage(page)
		if err != nil {
			return Page{}, fmt.Errorf("failed to create subpage, %w", err)
		}
		subPages = append(subPages, sp)
	}

	return Page{
		PageNumber: data.PageNumber,
		PrevPage:   data.PrevPage,
		NextPage:   data.NextPage,
		SubPages:   subPages,
		FetchedAt:  time.Now(),
	}, nil
}

func newSubPage(sp subPage) (SubPage, error) {
	gifReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(sp.GifAsBase64))
	img, err := gif.Decode(gifReader)
	if err != nil {
		return SubPage{}, fmt.Errorf("failed to decode gif: %w", err)
	}

	grid, err := newGrid(sp.AltText, img)
	if err != nil {
		return SubPage{}, fmt.Errorf("failed to generate grid: %w", err)
	}

	return SubPage{
		Number: sp.SubPageNumber,
		Text:   sp.AltText,
		Grid:   grid,
		Image:  img,
	}, nil
}

func newGrid(text string, img image.Image) (Grid, error) {
	colorMap := colorMap{
		image: img,
		dx:    img.Bounds().Dx() / (PageWidth - 1), // image is missing last col
		dy:    img.Bounds().Dy() / PageHeight,
	}

	rows := strings.Split(text, "\n")

	// text contains extra rows, this takes it from 30 to PageHeight (25)
	rows = rows[2 : len(rows)-3]

	var grid Grid
	for y, row := range rows {
		var cols Row

		x := 0 // range row returns byte index, we need rune index
		for _, c := range row {
			cols = append(cols, colorMap.getColor(x, y, c))
			x++
		}

		grid = append(grid, cols)
	}

	// add empty spaces to end of row (should be 1 missing)
	for y := range grid {
		for x := len(grid[y]); x < PageWidth-1; x++ {
			grid[y] = append(grid[y], colorMap.getColor(x, y, ' '))
		}
		grid[y] = append(grid[y], Cell{Char: ' ', Bg: black})
	}

	return grid, nil
}

type colorMap struct {
	image image.Image
	dx    int
	dy    int
}

func (c *colorMap) getColor(x, y int, char rune) Cell {
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

	return Cell{Char: char, Bg: bg, Fg: fg}
}
