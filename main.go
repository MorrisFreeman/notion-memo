package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MorrisFreeman/notion-memo/editor"
	notion "github.com/MorrisFreeman/notion-memo/notion"
)

func main() {
	notionKey := os.Getenv("NOTION_KEY")
	notionDatabaseId := os.Getenv("NOTION_DATABASE_ID")
	notionInboxBlockId := os.Getenv("NOTION_INBOX_BLOCK_ID")

	var n = flag.String("n", "", "name")
	var f = flag.String("f", "", "file")

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

	if *f != "" {
		f, err := os.Open(*f)
		if err != nil {
			os.Exit(1)
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			os.Exit(1)
		}

		notion.CreateDatabasePage(notionKey, notionDatabaseId, f.Name(), string(b))
		os.Exit(0)
	}

	body1 := flag.Arg(0)
	if body1 != "" {
		notion.CreateBlock(notionKey, notionInboxBlockId, "\n"+body1)
		os.Exit(0)
	}

	body2, err := editor.ReadEditor("")
	if err != nil {
		fmt.Printf("failed open editor: %s\n", err)
		os.Exit(1)
	}
	head := strings.Split(string(body2), "\n")[0]
	notion.CreateDatabasePage(notionKey, notionDatabaseId, head, string(body2))
	os.Exit(0)
}
