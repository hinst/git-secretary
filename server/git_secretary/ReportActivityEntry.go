package git_secretary

import git_stories_api "github.com/hinst/git-stories-api"

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

func (me *ReportActivityEntry) ReadRepositoryLogEntry(source *git_stories_api.RepositoryLogEntry) {
	me.ChangesetCount += 1
	for _, parent := range source.Parents {
		for _, diff := range parent.DiffRows {
			me.ChangedFileCount += 1
			me.InsertionCount += MaxOfTwoInts(0, diff.InsertionCount)
			me.DeletionCount += MaxOfTwoInts(0, diff.DeletionCount)
		}
	}
	me.Points = me.GetPoints()
}

func (me *ReportActivityEntry) Add(entry *ReportActivityEntry) {
	me.Points += entry.Points
	me.ChangesetCount += entry.ChangesetCount
	me.ChangedFileCount += entry.ChangedFileCount
	me.InsertionCount += entry.InsertionCount
	me.DeletionCount += entry.DeletionCount
}
