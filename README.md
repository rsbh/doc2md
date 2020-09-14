# doc2md

Fetch google docs from drive and save them as markdown in local

- Supports Fetching Spreadsheets
- Break Pages according to Table of contents (Experimetal)
- Support CodeBlocks (Experimetal)

## Config

| Key              | Description                                                                                                | Default |
| ---------------- | ---------------------------------------------------------------------------------------------------------- | ------- |
| folderId         | Drive folder ID, can be copied from address bar, all docs and sheet files from this folder will be fetched | ""      |
| docIds           | List of indivisuals Docs Ids to be fetched                                                                 | []      |
| outDir           | Output Destination                                                                                         | "out    |
| breakDoc         | Add support to break pages as per Table of contents (Experimetal)                                          | false   |
| supportCodeBlock | Add support for CodeBlocks (Experimetal)                                                                   | false   |
