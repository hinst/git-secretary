package git_client

import (
	"fmt"
	"strings"
)

type LogEntryRow struct {
	CommitHash   string   `json:"commitHash"`
	ParentHashes []string `json:"parentHashes"`
}

var _ fmt.Stringer = &LogEntryRow{}

func (this *LogEntryRow) String() string {
	return this.CommitHash + " <- " + strings.Join(this.ParentHashes, ",")
}

func (this *LogEntryRow) ParseGitLine(line string) {
	var parts = strings.Split(line, " ")
	this.CommitHash = parts[0]
	this.ParentHashes = nil
	for i := 1; i < len(parts); i++ {
		this.ParentHashes = append(this.ParentHashes, parts[i])
	}
}
