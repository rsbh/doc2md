package transformer

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

func getFileName(fullURLFile string, Ctype string) string {
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		panic(err)
	}
	ext, err := mime.ExtensionsByType(Ctype)
	if err != nil {
		panic(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	outputFile := fmt.Sprintf("%v%v", fileName, ext[0])
	return outputFile
}

// ReplaceImage fetch image from the link
func ReplaceImage(fullURLFile string) (string, []byte) {
	client := http.Client{}
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
	fileName := getFileName(fullURLFile, Ctype)
	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fileName, image
}
