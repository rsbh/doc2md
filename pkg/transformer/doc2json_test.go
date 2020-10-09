package transformer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/docs/v1"
)

func TestCleanText(t *testing.T) {
	t.Run("should replace new-line with <br/> if ignoreLineBreak is true", func(t *testing.T) {
		got := cleanText("hello world\n foo bar", true)
		want := "hello world<br/> foo bar"
		assert.Equal(t, want, got)
	})

	t.Run("should remove new line if ignoreLineBreak is false", func(t *testing.T) {
		got := cleanText("hello world\n foo bar", false)
		want := "hello world foo bar"
		assert.Equal(t, want, got)
	})

	t.Run("should trim the white space", func(t *testing.T) {
		got := cleanText("  hello world\n foo bar  ", false)
		want := "hello world foo bar"
		assert.Equal(t, want, got)
	})
}

func TestReplaceTags(t *testing.T) {
	t.Run("should replace html tag with `&lt;` and `&gt;`", func(t *testing.T) {
		text := "<hello/> foo bar"
		got := replaceTags(text)
		want := "&lt;hello/&gt; foo bar"
		assert.Equal(t, want, got)
	})

	t.Run("should not replace <br/> tag", func(t *testing.T) {
		text := "<br/> foo bar"
		got := replaceTags(text)
		want := "<br/> foo bar"
		assert.Equal(t, want, got)
	})
}

func TestGetText(t *testing.T) {
	t.Run("should return empty string if empty string is passed", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "",
				TextStyle: &docs.TextStyle{
					Bold: false,
				},
			},
		}
		got := getText(p, true, true)
		want := ""
		assert.Equal(t, want, got)
	})

	t.Run("should  format if text is italic", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Italic: true,
				},
			},
		}
		got := getText(p, true, false)
		want := "_Hello World_"
		assert.Equal(t, want, got)
	})

	t.Run("should format if text is bold", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Bold: true,
				},
			},
		}
		got := getText(p, true, false)
		want := "**Hello World**"
		assert.Equal(t, want, got)
	})

	t.Run("should not format if text is header", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Bold: true,
				},
			},
		}
		got := getText(p, true, true)
		want := "Hello World\n"
		assert.Equal(t, want, got)
	})

	t.Run("should format if text is both bold and italic", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Bold:   true,
					Italic: true,
				},
			},
		}
		got := getText(p, true, false)
		want := "**_Hello World_**"
		assert.Equal(t, want, got)
	})

	t.Run("should format if text is strikethrough", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Strikethrough: true,
				},
			},
		}
		got := getText(p, true, false)
		want := "~~Hello World~~"
		assert.Equal(t, want, got)
	})

	t.Run("should format if text is Link", func(t *testing.T) {
		p := &docs.ParagraphElement{
			TextRun: &docs.TextRun{
				Content: "Hello World",
				TextStyle: &docs.TextStyle{
					Link: &docs.Link{
						Url: "http://example.com",
					},
				},
			},
		}
		got := getText(p, true, false)
		want := "[Hello World](http://example.com)"
		assert.Equal(t, want, got)
	})
}

func TestGetListType(t *testing.T) {
	t.Run("should return ol by default", func(t *testing.T) {
		list := map[string]docs.List{}
		got := getListType(list, "")
		want := "ol"
		assert.Equal(t, want, got)
	})

	t.Run("should return ol if NestingLevel is empty", func(t *testing.T) {
		list := map[string]docs.List{
			"abc": docs.List{
				ListProperties: &docs.ListProperties{
					NestingLevels: []*docs.NestingLevel{},
				},
			},
		}
		got := getListType(list, "abc")
		want := "ol"
		assert.Equal(t, want, got)
	})

	t.Run("should return ol if GyphType is not empty", func(t *testing.T) {
		list := map[string]docs.List{
			"abc": docs.List{
				ListProperties: &docs.ListProperties{
					NestingLevels: []*docs.NestingLevel{
						&docs.NestingLevel{
							GlyphType: "DECIMAL",
						},
					},
				},
			},
		}
		got := getListType(list, "abc")
		want := "ol"
		assert.Equal(t, want, got)
	})

	t.Run("should return ul if GyphType is empty", func(t *testing.T) {
		list := map[string]docs.List{
			"abc": docs.List{
				ListProperties: &docs.ListProperties{
					NestingLevels: []*docs.NestingLevel{
						&docs.NestingLevel{
							GlyphType: "",
						},
					},
				},
			},
		}
		got := getListType(list, "abc")
		want := "ul"
		assert.Equal(t, want, got)
	})

	t.Run("should find the list id from list", func(t *testing.T) {
		list := map[string]docs.List{
			"abc": docs.List{
				ListProperties: &docs.ListProperties{
					NestingLevels: []*docs.NestingLevel{
						&docs.NestingLevel{
							GlyphType: "",
						},
					},
				},
			},
			"xyz": docs.List{
				ListProperties: &docs.ListProperties{
					NestingLevels: []*docs.NestingLevel{
						&docs.NestingLevel{
							GlyphType: "ROMAN",
						},
					},
				},
			},
		}
		got := getListType(list, "xyz")
		want := "ol"
		assert.Equal(t, want, got)
	})
}
