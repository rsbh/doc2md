package gdrive

import (
	"encoding/json"
	"fmt"
	"log"
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
func (s *Service) GetFiles(folderID string) {
	query := generateQuery(folderID)
	r, err := s.drive.Files.List().SupportsTeamDrives(true).IncludeTeamDriveItems(true).Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	if len(r.Files) == 0 {
		fmt.Println("No Drives found.")
	} else {
		for _, i := range r.Files {
			if i.MimeType == mimeTypeDocument {
				s.fetchDoc(i.Id)
			} else {
				fmt.Printf("%s (%s) (%s)\n", i.Name, i.Id, i.MimeType)
			}
		}
	}
}

func (s *Service) fetchDoc(docID string) {
	doc, err := s.doc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve doc: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(doc.Body, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
