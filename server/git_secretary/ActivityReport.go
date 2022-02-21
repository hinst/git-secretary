package git_secretary

import git_stories_api "github.com/hinst/git-stories-api"

type ActivityReport struct {
	Points           int `json:"points"`
	ChangesetCount   int `json:"changesetCount"`
	ChangedFileCount int `json:"changedFileCount"`
	InsertionCount   int `json:"insertionCount"`
	DeletionCount    int `json:"deletionCount"`
}

func (me *ActivityReport) GetPoints() int {
	return me.ChangesetCount + me.ChangedFileCount + me.InsertionCount + me.DeletionCount
}

func (me *ActivityReport) ReadRepositoryLogEntry(source *git_stories_api.RepositoryLogEntry) {
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

func (me *ActivityReport) Add(activityReport *ActivityReport) {
	me.Points += activityReport.Points
	me.ChangesetCount += activityReport.ChangesetCount
	me.ChangedFileCount += activityReport.ChangedFileCount
	me.InsertionCount += activityReport.InsertionCount
	me.DeletionCount += activityReport.DeletionCount
}
