package main

import (
	"log"
	"os"

	"github.com/aKjeller/text-tv/internal/display"
	"github.com/aKjeller/text-tv/internal/tui"
	"github.com/aKjeller/text-tv/pkg/svttext"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	if len(os.Args) == 2 {
		simple()
	} else {
		interactive()
	}
}

func simple() {
	pageId := os.Args[1]

	page, err := svttext.GetPage(pageId)
	if err != nil {
		log.Fatalf("failed to get page: %v", err)
	}

	display.RenderPage(page)
}

func interactive() {
	p := tea.NewProgram(
		tui.Model{},
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
