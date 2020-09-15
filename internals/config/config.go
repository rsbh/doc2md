package config

// Configurations exported
type Configurations struct {
	FolderID         string   `json:"folderId" yaml:"folderId"`
	DocIDs           []string `json:"docIds" yaml:"docIds"`
	OutDir           string   `json:"outDir" yaml:"outDir"`
	BreakDoc         bool     `json:"breakDoc" yaml:"breakDoc"`
	SupportCodeBlock bool     `json:"supportCodeBlock" yaml:"supportCodeBlock"`
	ExtendendQuery   bool     `json:"extendedQuery" yaml:"extendedQuery"`
}
