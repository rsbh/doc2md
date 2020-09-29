package transformer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestReplaceImage(t *testing.T) {
	t.Run("should return the image name after fetching image", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution
		src := "http://foo.com/img"
		xbody := []byte{'I', 'M', 'A', 'G', 'E'}
		body := bytes.NewBuffer(xbody)

		gock.New(src).
			Reply(200).
			Body(body).
			SetHeader("Content-Type", "image/png")

		fileName, image := ReplaceImage(src)
		want := "img.png"
		assert.Equal(t, want, fileName)
		assert.Equal(t, xbody, image)
	})
}

func TestGetFileName(t *testing.T) {
	t.Run("should return filename from url", func(t *testing.T) {
		url := "http://example.com/a"
		got := getFileName(url, "image/png")
		want := "a.png"
		assert.Equal(t, want, got)
	})

	t.Run("should return filename from nexted url", func(t *testing.T) {
		url := "http://example.com/a/b"
		got := getFileName(url, "image/png")
		want := "b.png"
		assert.Equal(t, want, got)
	})

	t.Run("should return change extension by Content Type", func(t *testing.T) {
		url := "http://example.com/a/b"
		got := getFileName(url, "image/jpeg")
		want := "b.jpg"
		assert.Equal(t, want, got)
	})
}
