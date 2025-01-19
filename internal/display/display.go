package display

import (
	"fmt"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
)

func RenderPage(page svttext.Page) error {
	g := createGrid(page.Text)
	g = g[2 : len(g)-3]

	colors, err := newColorMap(page.Image, g.getWidth(), len(g))
	if err != nil {
		return err
	}

	g.render(colors)

	return nil
}

type grid [][]rune

func (g grid) getWidth() int {
	width := 0
	for _, r := range g {
		width = max(width, len(r))
	}
	return width
}

func (g grid) render(colors *colorMap) {
	for i, row := range g {
		r := ""
		for j, c := range row {
			r += colors.getColor(j, i).colorRune(c)
		}
		r += colors.getColor(i, len(row)-1).colorRune(' ')
		fmt.Println(r)
	}
}

func createGrid(text string) grid {
	var grid grid
	rows := strings.Split(text, "\n")
	for _, row := range rows {
		var cols []rune
		for _, c := range row {
			cols = append(cols, c)
		}
		grid = append(grid, cols)
	}

	width := grid.getWidth()

	for i := range grid {
		for j := len(grid[i]); j < width; j++ {
			grid[i] = append(grid[i], ' ')
		}
	}

	return grid
}
