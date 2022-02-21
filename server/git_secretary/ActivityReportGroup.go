package git_secretary

import "time"

type ActivityReportGroup struct {
	Time     time.Time      `json:"time"`
	Period   time.Duration  `json:"duration"`
	Activity ActivityReport `json:"activity"`
	// Key is author name
	Authors map[string]*ActivityReport `json:"authors"`
}

type ActivityReportGroups []*ActivityReportGroup

func (me *ActivityReportGroup) Aggregate(authorName string, activityReport *ActivityReport) {
	me.Activity.Add(activityReport)
	var authorActivityEntry = me.getOrCreateAuthor(authorName)
	authorActivityEntry.Add(activityReport)
}

func (me *ActivityReportGroup) getOrCreateAuthor(authorName string) (activityEntry *ActivityReport) {
	if nil == me.Authors {
		me.Authors = make(map[string]*ActivityReport)
	}
	activityEntry = me.Authors[authorName]
	if nil == activityEntry {
		activityEntry = &ActivityReport{}
		me.Authors[authorName] = activityEntry
	}
	return activityEntry
}
