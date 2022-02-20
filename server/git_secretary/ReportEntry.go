package git_secretary

import "time"

type ReportEntry struct {
	Time     time.Time           `json:"time"`
	Period   time.Duration       `json:"duration"`
	Activity ReportActivityEntry `json:"activity"`
	// Key is author name
	Authors map[string]*ReportActivityEntry `json:"authors"`
}

type ReportEntries []*ReportEntry

func (me *ReportEntry) Aggregate(reportEntry *ReportEntry) {
	for authorName, activityEntry := range reportEntry.Authors {
		me.Activity.Add(activityEntry)
		var authorActivityEntry = me.getOrCreateAuthor(authorName)
		authorActivityEntry.Add(activityEntry)
	}
}

func (me *ReportEntry) getOrCreateAuthor(authorName string) (activityEntry *ReportActivityEntry) {
	if nil == me.Authors {
		me.Authors = make(map[string]*ReportActivityEntry)
	}
	activityEntry = me.Authors[authorName]
	if nil == activityEntry {
		activityEntry = &ReportActivityEntry{}
		me.Authors[authorName] = activityEntry
	}
	return activityEntry
}
