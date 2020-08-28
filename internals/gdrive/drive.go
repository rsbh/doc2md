package gdrive

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

const mimeTypeDocument = "application/vnd.google-apps.document"
const mimeTypeFolder = "application/vnd.google-apps.folder"
const mimeTypeSheet = "application/vnd.google-apps.spreadsheet"

// Service to fetch files
type Service struct {
	driveSrv *drive.Service
	docSrv   *docs.Service
}

// GetFiles return google drive files
func (s Service) GetFiles(folderID string) {
	var q = ""
	if folderID != "" {
		q = fmt.Sprintf("'%v' in parents and ", folderID)
	}
	query := fmt.Sprintf("%v(mimeType='%v' or mimeType='%v' or mimeType='%v') and trashed = false", q, mimeTypeDocument, mimeTypeFolder, mimeTypeSheet)
	fmt.Println(query)
}

// Init initialize the services
func (s Service) Init(client *http.Client) {
	driveSrv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	docsSrv, err := docs.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Doc client: %v", err)
	}
	s.driveSrv = driveSrv
	s.docSrv = docsSrv
}
