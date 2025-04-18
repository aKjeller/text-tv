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

func PrintPage(pageId string) {
	client := svttext.NewClient()
	page, err := client.GetPage(pageId)
	if err != nil {
		log.Fatalf("failed to get page: %v", err)
	}

	fmt.Println()
	for chunk := range slices.Chunk(page.SubPages, getDisplayWidth(svttext.PageWidth)) {
		for row := range svttext.PageHeight {
			var sb strings.Builder
			for _, sp := range chunk {
				sb.WriteString(sp.Grid[row].ColorString())
			}
			fmt.Println(sb.String())
		}
	}
}

func getDisplayWidth(pageWidth int) int {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 1
	}

	width, _, err := term.GetSize(fd)
	if err != nil {
		log.Fatalf("could not get terminal size: %v", err)
	}

	return width / pageWidth
}
