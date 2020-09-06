package gdrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/rsbh/doc2md/pkg/transformer"
	"github.com/spf13/viper"
)

func (s *Service) fetchDoc(docID string, bc []string) {
	outDir := viper.GetString("OutDir")
	breadCrumbs := path.Join(bc...)
	doc, err := s.doc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve doc: %v", err)
	}

	outPath := path.Join(outDir, breadCrumbs, doc.Title)

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		os.MkdirAll(outPath, 0700) // Create your file
	}

	_, _, toc, pages := transformer.DocToJSON(doc, true)

	for _, p := range pages {
		x := transformer.JSONToMD(p.Contents)
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
