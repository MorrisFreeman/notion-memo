package notion

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dstotijn/go-notion"
)

type httpTransport struct {
	w io.Writer
}

func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}

func buildText(text string) *notion.Text {
	return &notion.Text{
		Content: text,
	}
}

func buildRichText(text string) []notion.RichText {
	return []notion.RichText{
		{
			Text: buildText(text),
		},
	}
}

func parseLine(line string) notion.Block {
	if strings.HasPrefix(line, "# ") {
		l := strings.Replace(line, "# ", "", 1)
		return notion.Heading1Block{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "## ") {
		l := strings.Replace(line, "## ", "", 1)
		return notion.Heading2Block{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "### ") {
		l := strings.Replace(line, "### ", "", 1)
		return notion.Heading3Block{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "- ") {
		l := strings.Replace(line, "- ", "", 1)
		return notion.BulletedListItemBlock{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "* ") {
		l := strings.Replace(line, "* ", "", 1)
		return notion.BulletedListItemBlock{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "1. ") {
		l := strings.Replace(line, "1. ", "", 1)
		return notion.NumberedListItemBlock{
			RichText: buildRichText(l),
		}
	}

	if strings.HasPrefix(line, "-[] ") {
		l := strings.Replace(line, "-[] ", "", 1)
		return notion.ToDoBlock{
			RichText: buildRichText(l),
		}
	}

	return notion.ParagraphBlock{
		RichText: buildRichText(line),
	}

}

func parseBody(body string) []notion.Block {
	scanner := bufio.NewScanner(strings.NewReader(body))
	var blocks []notion.Block
	for scanner.Scan() {
		line := parseLine(scanner.Text())
		blocks = append(blocks, line)
	}

	return blocks
}

func buildCreatePageParams(databaseId, name, body string) (*notion.CreatePageParams, error) {
	p := &notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   databaseId,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: name,
						},
					},
				},
			},
		},
	}
	if body != "" {
		children := parseBody(body)
		p.Children = children
	}

	return p, nil
}

func CreateDatabasePage(notionKey, notionDatabaseId, name, body string) {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(notionKey, notion.WithHTTPClient(httpClient))
	params, _ := buildCreatePageParams(notionDatabaseId, name, body)

	_, err := client.CreatePage(ctx, *params)
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		log.Fatal(err)
	}
}

func CreateBlock(notionKey, blockId, text string) {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(notionKey, notion.WithHTTPClient(httpClient))

	fetched := RetrieveBlock(notionKey, blockId)
	pb, _ := fetched.(*notion.ParagraphBlock)
	r := buildRichText(text)
	pb.RichText = append(pb.RichText, r[0])

	_, err := client.UpdateBlock(ctx, blockId, pb)
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		log.Fatal(err)
	}
}

func RetrieveBlock(notionKey, blockId string) notion.Block {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(notionKey, notion.WithHTTPClient(httpClient))

	block, _ := client.FindBlockByID(ctx, blockId)

	return block
}
