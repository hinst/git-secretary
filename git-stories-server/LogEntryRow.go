package main

import (
	"fmt"
	"strings"
)

type LogEntryRow struct {
	CommitHash   string   `json:"commitHash"`
	ParentHashes []string `json:"parentHashes"`
}

var _ fmt.Stringer = &LogEntryRow{}

func (row *LogEntryRow) String() string {
	return row.CommitHash + " <- " + strings.Join(row.ParentHashes, ",")
}

func (row *LogEntryRow) ParseGitLine(line string) {
	var parts = strings.Split(line, " ")
	row.CommitHash = parts[0]
	row.ParentHashes = nil
	for i := 1; i < len(parts); i++ {
		row.ParentHashes = append(row.ParentHashes, parts[i])
	}
}
