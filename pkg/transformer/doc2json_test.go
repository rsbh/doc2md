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
		assert.Equal(t, got, want)
	})

	t.Run("should remove new line if ignoreLineBreak is false", func(t *testing.T) {
		got := cleanText("hello world\n foo bar", false)
		want := "hello world foo bar"
		assert.Equal(t, got, want)
	})

	t.Run("should trim the white space", func(t *testing.T) {
		got := cleanText("  hello world\n foo bar  ", false)
		want := "hello world foo bar"
		assert.Equal(t, got, want)
	})
}

func textReplaceTags(t *testing.T) {
	t.Run("should replace html tag with `&lt;` and `&gt;`", func(t *testing.T) {
		text := "<hello/> foo bar"
		got := replaceTags(text)
		want := "&lt;hello/&gt; foo bar\n"
		assert.Equal(t, got, want)
	})

	t.Run("should not replace <br/> tag", func(t *testing.T) {
		text := "<br/> foo bar"
		got := replaceTags(text)
		want := "&lt;hello/&gt; foo bar\n"
		assert.Equal(t, got, want)
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
		assert.Equal(t, got, want)
	})
}
