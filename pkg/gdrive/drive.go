package gdrive

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"google.golang.org/api/drive/v3"
)

func generateQuery(folderID string) string {
	extendedQuery := viper.GetString("extendedQuery")
	var f = ""
	if folderID != "" {
		f = fmt.Sprintf("'%v' in parents and ", folderID)
	}
	query := fmt.Sprintf("%v(mimeType='%v' or mimeType='%v' or mimeType='%v') and trashed = false %v", f, mimeTypeDocument, mimeTypeFolder, mimeTypeSheet, extendedQuery)
	return query
}

func (s *Service) getFiles(query string) []*drive.File {
	r, err := s.drive.Files.List().SupportsTeamDrives(true).IncludeTeamDriveItems(true).Q(query).Fields("files(id, name, mimeType, description, createdTime, modifiedTime, lastModifyingUser)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	return r.Files
}

// GetFiles return google drive files
func (s *Service) GetFiles(folderID string, bc []string, rwg *sync.WaitGroup) {
	var wg sync.WaitGroup
	query := generateQuery(folderID)
	files := s.getFiles(query)
	if len(files) == 0 {
		fmt.Println("No Files found. FolderID : ", folderID)
	} else {
		for _, i := range files {
			if i.MimeType == mimeTypeDocument {
				wg.Add(1)
				meta := FrontMatter{"", i.Description, i.LastModifyingUser.DisplayName, i.ModifiedTime, i.CreatedTime}
				go s.FetchDoc(i.Id, bc, meta, &wg)
			} else if i.MimeType == mimeTypeFolder {
				wg.Add(1)
				newBc := append(bc, i.Name)
				go s.GetFiles(i.Id, newBc, &wg)
			} else if i.MimeType == mimeTypeSheet {
				wg.Add(1)
				go s.fetchSheet(i.Id, i.Name, bc, &wg)
			} else {
				fmt.Printf("%s (%s) (%s)\n", i.Name, i.Id, i.MimeType)
			}
		}
	}
	wg.Wait()
	rwg.Done()
}
