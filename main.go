package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MorrisFreeman/notion-memo/editor"
	notion "github.com/MorrisFreeman/notion-memo/notion"
)

func main() {
	var i = flag.String("i", "", "instantly")
	flag.Parse()

	if *i != "" {
		notion.Create(*i, "")
		os.Exit(0)
	}

	name := "test" // flag.Arg(0)
	if name != "" {
		body, err := editor.ReadEditor()
		fmt.Printf("%s\n", body)
		if err != nil {
			fmt.Printf("failed open editor: %s\n", err)
			os.Exit(1)
		}
		notion.Create(name, string(body))
	}
}
