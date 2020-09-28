package transformer

import (
	"testing"
)

func TestCleanText(t *testing.T) {
	t.Run("should replace new-line with <br/> if ignoreLineBreak is true", func(t *testing.T) {
		got := cleanText("hello world\n foo bar", true)
		want := "hello world<br/> foo bar"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("should remove new line if ignoreLineBreak is false", func(t *testing.T) {
		got := cleanText("hello world\n foo bar", false)
		want := "hello world foo bar"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("should trim the white space", func(t *testing.T) {
		got := cleanText("  hello world\n foo bar  ", false)
		want := "hello world foo bar"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
