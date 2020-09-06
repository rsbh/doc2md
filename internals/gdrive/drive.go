package gdrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

func generateQuery(folderID string) string {
	var f = ""
	if folderID != "" {
		f = fmt.Sprintf("'%v' in parents and ", folderID)
	}
	query := fmt.Sprintf("%v(mimeType='%v' or mimeType='%v' or mimeType='%v') and trashed = false", f, mimeTypeDocument, mimeTypeFolder, mimeTypeSheet)
	return query
}

// GetFiles return google drive files
func (s *Service) GetFiles(folderID string, bc []string) {
	query := generateQuery(folderID)
	r, err := s.drive.Files.List().SupportsTeamDrives(true).IncludeTeamDriveItems(true).Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	if len(r.Files) == 0 {
		fmt.Println("No Drives found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s) (%s)\n", i.Name, i.Id, i.MimeType)
			if i.MimeType == mimeTypeDocument {
				s.fetchDoc(i.Id, bc)
			} else if i.MimeType == mimeTypeFolder {
				newBc := append(bc, i.Name)
				s.GetFiles(i.Id, newBc)
			} else {
				fmt.Printf("%s (%s) (%s)\n", i.Name, i.Id, i.MimeType)
			}
		}
	}
}

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

	outputFile := path.Join(outPath, "index.json")

	prettyJSON, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	_ = ioutil.WriteFile(outputFile, prettyJSON, 0644)
}
