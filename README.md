# doc2md

![GitHub Workflow Status](https://github.com/rsbh/doc2md/workflows/CI/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/rsbh/doc2md/badge.svg)](https://coveralls.io/github/rsbh/doc2md) [![Go Report Card](https://goreportcard.com/badge/github.com/rsbh/doc2md)](https://goreportcard.com/report/github.com/rsbh/doc2md)

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
| extendedQuery    | Need to extend the drive list files query                                                                  | ""      |
