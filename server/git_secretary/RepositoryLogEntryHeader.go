package git_secretary

import (
	"fmt"
	"strings"
)

type RepositoryLogEntryHeader struct {
	CommitHash   string   `json:"commitHash"`
	ParentHashes []string `json:"parentHashes"`
}

var _ fmt.Stringer = &RepositoryLogEntryHeader{}

func (row *RepositoryLogEntryHeader) String() string {
	return row.CommitHash + " <- " + strings.Join(row.ParentHashes, ",")
}

func (row *RepositoryLogEntryHeader) ParseGitLine(line string) {
	var parts = strings.Split(line, " ")
	row.CommitHash = parts[0]
	row.ParentHashes = nil
	for i := 1; i < len(parts); i++ {
		row.ParentHashes = append(row.ParentHashes, parts[i])
	}
}

type RepositoryLogEntryHeaders []RepositoryLogEntryHeader

func (rows RepositoryLogEntryHeaders) GetPortions(portionSize int) (portions []RepositoryLogEntryHeaders) {
	for _, row := range rows {
		if len(portions) == 0 || len(portions[len(portions)-1]) >= portionSize {
			portions = append(portions, nil)
		}
		portions[len(portions)-1] = append(portions[len(portions)-1], row)
	}
	return
}
