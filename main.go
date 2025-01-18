package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	text := page.Text
	text = strings.TrimPrefix(text, "\n")
	text = strings.TrimSuffix(text, "\n\n\n")

	fmt.Println(text)
}
