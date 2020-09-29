package transformer

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

// ReplaceImage fetch image from the link
func ReplaceImage(fullURLFile string) (string, []byte) {
	client := http.Client{}

	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		panic(err)
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	req, err := http.NewRequest("GET", fullURLFile, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	Ctype := resp.Header.Get("Content-Type")

	ext, err := mime.ExtensionsByType(Ctype)
	if err != nil {
		panic(err)
	}

	outputFile := fmt.Sprintf("%v%v", fileName, ext[0])
	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return outputFile, image
}
