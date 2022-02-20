package git_secretary

import "time"

type ReportEntry struct {
	Time     time.Time                      `json:"time"`
	Period   time.Duration                  `json:"duration"`
	Activity ReportActivityEntry            `json:"activity"`
	Authors  map[string]ReportActivityEntry `json:"authors"`
}

type ReportEntries []*ReportEntry

type ReportActivityEntry struct {
	Points           int `json:"points"`
	ChangesetCount   int `json:"changesetCount"`
	ChangedFileCount int `json:"changedFileCount"`
	InsertionCount   int `json:"insertionCount"`
	DeletionCount    int `json:"deletionCount"`
}
