package git_secretary

import (
	"time"

	git_stories_api "github.com/hinst/git-stories-api"
)

type ReportEntry struct {
	Time     time.Time                      `json:"time"`
	Period   time.Duration                  `json:"duration"`
	Activity ReportActivityEntry            `json:"activity"`
	Authors  map[string]ReportActivityEntry `json:"authors"`
}

type ReportEntries []*ReportEntry

func (me *ReportEntry) Aggregate(entry *ReportEntry) {

}

type ReportActivityEntry struct {
	Points           int `json:"points"`
	ChangesetCount   int `json:"changesetCount"`
	ChangedFileCount int `json:"changedFileCount"`
	InsertionCount   int `json:"insertionCount"`
	DeletionCount    int `json:"deletionCount"`
}

func (me *ReportActivityEntry) GetPoints() int {
	return me.ChangesetCount + me.ChangedFileCount + me.InsertionCount + me.DeletionCount
}

func (me *ReportActivityEntry) readRepositoryLogEntry(source *git_stories_api.RepositoryLogEntry) {
	me.ChangesetCount += 1
	for _, parent := range source.Parents {
		me.ChangedFileCount += len(parent.DiffRows)
		for _, diff := range parent.DiffRows {
			me.ChangedFileCount += 1
			me.InsertionCount += MaxOfTwoInts(0, diff.InsertionCount)
			me.DeletionCount += MaxOfTwoInts(0, diff.DeletionCount)
		}
	}
	me.Points = me.GetPoints()
}
