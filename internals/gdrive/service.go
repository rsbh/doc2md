package gdrive

import (
	"log"
	"net/http"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

const mimeTypeDocument = "application/vnd.google-apps.document"
const mimeTypeFolder = "application/vnd.google-apps.folder"
const mimeTypeSheet = "application/vnd.google-apps.spreadsheet"

// Service to fetch files
type Service struct {
	drive  *drive.Service
	doc    *docs.Service
	sheets *sheets.Service
}

// Init initialize the services
func (s *Service) Init(client *http.Client) {
	driveSrv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	docsSrv, err := docs.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Doc client: %v", err)
	}
	sheetsSrv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	s.drive = driveSrv
	s.doc = docsSrv
	s.sheets = sheetsSrv
}
