package git_client

import "time"

type DetailedLogEntryRow struct {
	Time       time.Time
	Parents    []ParentInfoEntry
	CommitHash string
}

type ParentInfoEntry struct {
	CommitHash  string
	DiffSummary []DiffSummaryRow
}
