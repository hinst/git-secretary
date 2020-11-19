package git_client

import "time"

type DetailedLogEntryRow struct {
	LogEntry    LogEntryRow
	Time        time.Time
	DiffSummary []DiffSummaryRow
}
