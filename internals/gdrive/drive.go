package gdrive

import (
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
			if i.MimeType == mimeTypeDocument {
				s.fetchDoc(i.Id, bc)
			} else if i.MimeType == mimeTypeFolder {
				newBc := append(bc, i.Name)
				s.GetFiles(i.Id, newBc)
			} else if i.MimeType == mimeTypeSheet {
				s.fetchSheet(i.Id, i.Name, bc)
			} else {
				fmt.Printf("%s (%s) (%s)\n", i.Name, i.Id, i.MimeType)
			}
		}
	}
}
