package display

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
	"golang.org/x/term"
)

func RenderPage(page svttext.Page) error {
	var pages []grid
	for _, sp := range page.SubPages {
		g := createGrid(sp.Text)
		g = g[2 : len(g)-3]
		g = toColorGrid(g, sp)
		pages = append(pages, g)
	}

	for chunk := range slices.Chunk(pages, getDisplayWidth(len(pages[0][0]))) {
		for row, _ := range chunk[0] {
			r := ""
			for _, p := range chunk {
				r += p.encodeRow(row)
			}
			fmt.Println(r)
		}
	}

	return nil
}

type grid [][]cell

func (g grid) encodeRow(index int) string {
	r := ""
	for _, c := range g[index] {
		r += c.colorRune()
	}
	return r
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

func getDisplayWidth(pageWidth int) int {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 1
	}

	width, _, err := term.GetSize(fd)
	if err != nil {
		log.Fatalf("could not get terminal size", err)
	}

	return width / pageWidth
}
