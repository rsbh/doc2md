package gdrive

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/spf13/viper"
)

func sheetToJSON(values [][]interface{}) []map[string]string {
	var data []map[string]string
	keys, body := values[0], values[1:]
	for _, row := range body {
		var rowData = make(map[string]string)
		for i, k := range keys {
			key := fmt.Sprintf("%v", k)
			value := ""
			if len(row) > 0 && len(row) > i {
				value = fmt.Sprintf("%v", row[i])
			}
			rowData[key] = value
		}
		data = append(data, rowData)
	}
	return data
}

func (s *Service) fetchSheet(spreadsheetId string, name string, bc []string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := s.sheets.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var ranges []string

	for _, i := range resp.Sheets {
		ranges = append(ranges, i.Properties.Title)
	}

	sheet, err := s.sheets.Spreadsheets.Values.BatchGet(spreadsheetId).Ranges(ranges...).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	for i, value := range sheet.ValueRanges {
		data := sheetToJSON(value.Values)
		fileName := fmt.Sprintf("%v.json", ranges[i])
		saveSheet(data, name, fileName, bc)
	}

}

func saveSheet(data []map[string]string, folder string, fileName string, bc []string) {
	outDir := viper.GetString("OutDir")
	breadCrumbs := path.Join(bc...)
	outPath := path.Join(outDir, breadCrumbs, folder)
	json, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Unable to parse sheet data: %v", err)
	}
	c := FetchedDoc{outPath, fileName, json}
	c.SaveToFile()
}
