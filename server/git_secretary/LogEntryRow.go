package git_secretary

import (
	"fmt"
	"strings"
)

// TODO rename RepositoryLogEntryHeader
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

type LogEntryRows []LogEntryRow

func (rows LogEntryRows) GetPortions(portionSize int) (portions []LogEntryRows) {
	for _, row := range rows {
		if len(portions) == 0 || len(portions[len(portions)-1]) >= portionSize {
			portions = append(portions, nil)
		}
		portions[len(portions)-1] = append(portions[len(portions)-1], row)
	}
	return
}
