package notion

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dstotijn/go-notion"
)

func TestBuildText(t *testing.T) {
	tests := []struct {
		text     string
		expected *notion.Text
	}{
		{
			"",
			&notion.Text{
				Content: "",
			},
		},
		{
			"test",
			&notion.Text{
				Content: "test",
			},
		},
	}

	for _, tt := range tests {
		txt := buildText(tt.text)
		if *tt.expected != *txt {
			ex, _ := json.Marshal(tt.expected)
			got, _ := json.Marshal(txt)
			t.Errorf("txt is not %s, got=%s", string(ex), string(got))
		}
	}
}

func TestBuildRitchText(t *testing.T) {
	tests := []struct {
		text     string
		expected []notion.RichText
	}{
		{
			"",
			[]notion.RichText{
				{
					Text: &notion.Text{
						Content: "",
					},
				},
			},
		},
		{
			"test",
			[]notion.RichText{
				{
					Text: &notion.Text{
						Content: "test",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		rt := buildRichText(tt.text)
		if len(rt) != 1 {
			t.Errorf("rt is wrong num of elements, got=%d", len(rt))
			continue
		}
		if !reflect.DeepEqual(rt[0], tt.expected[0]) {
			ex, _ := json.Marshal(tt.expected[0])
			got, _ := json.Marshal(rt[0])
			t.Errorf("rt is not %s, got=%s", string(ex), string(got))
			continue
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		text     string
		expected []notion.RichText
	}{
		{
			"",
			[]notion.RichText{
				{
					Text: &notion.Text{
						Content: "",
					},
				},
			},
		},
		{
			"test",
			[]notion.RichText{
				{
					Text: &notion.Text{
						Content: "test",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		rt := buildRichText(tt.text)
		if len(rt) != 1 {
			t.Errorf("rt is wrong num of elements, got=%d", len(rt))
			continue
		}
		if !reflect.DeepEqual(rt[0], tt.expected[0]) {
			ex, _ := json.Marshal(tt.expected[0])
			got, _ := json.Marshal(rt[0])
			t.Errorf("rt is not %s, got=%s", string(ex), string(got))
			continue
		}
	}
}
