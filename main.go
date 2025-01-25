package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aKjeller/text-tv/internal/display"
	"github.com/aKjeller/text-tv/pkg/svttext"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <page_id>")
		os.Exit(1)
	}

	pageId := os.Args[1]

	page, err := svttext.GetPage(pageId)
	if err != nil {
		log.Fatalf("failed to get page", err)
	}

	err = display.RenderPage(page)
	if err != nil {
		log.Fatalf("failed to render page", err)
	}
}
