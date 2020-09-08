package gdrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	t "github.com/rsbh/doc2md/pkg/transformer"
	"github.com/spf13/viper"
)

func replaceImages(p []t.TagContent, imageFolder string) []t.TagContent {

	for i, tc := range p {
		img := tc["img"].Image
		if img.Source != "" {
			name, content := t.ReplaceImage(img.Source)
			imgPath := path.Join(imageFolder, name)
			ioutil.WriteFile(imgPath, content, 0644)
			imgLink := path.Join("images", name)
			p[i] = t.TagContent{"img": {"", t.ImageObject{imgLink, img.Title, img.Description}, t.Table{}, t.CodeBlock{}}}
		}
	}
	return p
}

// FetchDoc fetch google doc from drive
func (s *Service) FetchDoc(docID string, bc []string) {
	outDir := viper.GetString("OutDir")
	breadCrumbs := path.Join(bc...)
	doc, err := s.doc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve doc: %v", err)
	}

	outPath := path.Join(outDir, breadCrumbs, doc.Title)
	imageFolder := path.Join(outPath, "images")

	if _, err := os.Stat(imageFolder); os.IsNotExist(err) {
		os.MkdirAll(imageFolder, 0700) // Create your file
	}

	pages, toc, _ := t.DocToJSON(doc, true)

	for _, p := range pages {
		updatedContent := replaceImages(p.Contents, imageFolder)
		x := t.JSONToMD(updatedContent)
		md := []byte(x)
		fileName := fmt.Sprintf("%v.md", p.Title)

		outputFile := path.Join(outPath, fileName)

		_ = ioutil.WriteFile(outputFile, md, 0644)

	}
	prettyToc, err := json.MarshalIndent(toc, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	tocFilePath := path.Join(outPath, "toc.json")

	_ = ioutil.WriteFile(tocFilePath, prettyToc, 0644)

}
