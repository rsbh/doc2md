package transformer

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func ReplaceImage(fullUrlFile string) (string, []byte) {
	client := httpClient()
	fileURL, err := url.Parse(fullUrlFile)
	path := fileURL.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	if err != nil {
		panic(err)
	}
	resp, err := client.Get(fullUrlFile)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	Ctype := resp.Header.Get("Content-Type")
	ext, _ := mime.ExtensionsByType(Ctype)
	outputFile := fmt.Sprintf("%v%v", fileName, ext[0])
	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return outputFile, image
}
