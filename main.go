package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MorrisFreeman/notion-memo/editor"
	notion "github.com/MorrisFreeman/notion-memo/notion"
)

func main() {
	notionKey := os.Getenv("NOTION_KEY")
	notionDatabaseId := os.Getenv("NOTION_DATABASE_ID")

	var n = flag.String("n", "", "name")
	flag.Parse()

	if *n != "" {
		body, err := editor.ReadEditor("")
		if err != nil {
			fmt.Printf("failed open editor: %s\n", err)
			os.Exit(1)
		}
		notion.CreateDatabasePage(notionKey, notionDatabaseId, *n, string(body))
		os.Exit(0)
	}

	name := flag.Arg(0)
	body, err := editor.ReadEditor(name)
	if len(body) == 0 {
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("failed open editor: %s\n", err)
		os.Exit(1)
	}
	notion.CreateDatabasePage(notionKey, notionDatabaseId, name, string(body))
	os.Exit(0)
}
