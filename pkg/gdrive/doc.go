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
	"gopkg.in/yaml.v2"
)

type FetchedDoc struct {
	OutPath  string
	FileName string
	Data     []byte
}

type FrontMatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"Description"`
	ModifiedBy  string `yaml:"modifiedBy"`
	ModifiedAt  string `yaml:"modifiedAt"`
	CreatedAt   string `yaml:"createdAt"`
}

// FetchDoc fetch google doc from drive
func (s *Service) FetchDoc(docID string, bc []string, meta FrontMatter) {
	outDir := viper.GetString("OutDir")
	breakDoc := viper.GetBool("BreakDoc")
	supportCodeBlock := viper.GetBool("SupportCodeBlock")

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

	pages, toc := t.DocToJSON(doc, imageFolder, supportCodeBlock, breakDoc)

	for _, p := range pages {
		meta.Title = p.Title
		frontMatter, err := yaml.Marshal(&meta)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fileName := fmt.Sprintf("%v.md", p.Title)

		md := t.JSONToMD(p.Contents)
		data := fmt.Sprintf("---\n%v\n---\n\n%v", string(frontMatter), md)

		content := []byte(data)
		d := FetchedDoc{outPath, fileName, content}
		d.SaveToFile()
	}
	prettyToc, err := json.MarshalIndent(toc, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	t := FetchedDoc{outPath, "toc.json", prettyToc}
	t.SaveToFile()
}

func (d FetchedDoc) SaveToFile() {
	if _, err := os.Stat(d.OutPath); os.IsNotExist(err) {
		os.MkdirAll(d.OutPath, 0700) // Create your file
	}

	outputFile := path.Join(d.OutPath, d.FileName)
	ioutil.WriteFile(outputFile, d.Data, 0644)
}
