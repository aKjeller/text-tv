package display

import (
	"fmt"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
)

func RenderPage(page svttext.Page) error {
	var pages []grid
	for _, sp := range page.SubPages {
		g := createGrid(sp.Text)
		g = g[2 : len(g)-3]
		g = toColorGrid(g, sp)
		pages = append(pages, g)
	}

	for _, p := range pages {
		p.render()
	}

	return nil
}

type grid [][]cell

func (g grid) getWidth() int {
	width := 0
	for _, r := range g {
		width = max(width, len(r))
	}
	return width
}

func (g grid) render() {
	for _, row := range g {
		r := ""
		for _, c := range row {
			r += c.colorRune()
		}

		// for non transparent background
		//r += cell{char: ' ', bg: color.Black}.colorRune()

		fmt.Println(r)
	}
}

func createGrid(text string) grid {
	var grid grid
	var width int

	rows := strings.Split(text, "\n")
	for _, row := range rows {
		var cols []cell
		for _, c := range row {
			cols = append(cols, cell{char: c})
		}
		width = max(width, len(cols))
		grid = append(grid, cols)
	}

	// add empty spaces to end of row
	for i := range grid {
		for j := len(grid[i]); j < width; j++ {
			grid[i] = append(grid[i], cell{char: ' '})
		}
	}

	return grid
}
