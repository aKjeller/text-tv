package main

import (
	"log"
	"os"

	"github.com/aKjeller/text-tv/internal/display"
	"github.com/aKjeller/text-tv/internal/tui"
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
	display.PrintPage(os.Args[1])
}

func interactive() {
	p := tea.NewProgram(
		tui.NewModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
