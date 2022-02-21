package git_secretary

import (
	"sort"
	"time"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
)

type ReportGenerator struct {
	Storage       *Storage
	Configuration *Configuration
	Request       ReadReportRequest
	Update        WebTaskUpdateFunc
}

func (me *ReportGenerator) Generate() {
	var gitClient = (&CachedGitClient{}).Create(me.Storage, me.Request.Directory)
	gitClient.SetProgressReceiver(func(total int, done int) {
		me.Update(&WebTask{Activity: "Read repository", Total: total, Done: done})
	})
	var repositoryLogEntries, gitError = gitClient.ReadDetailedLog(0)
	if nil != gitError {
		me.Update(&WebTask{Error: gitError.Error()})
		return
	}
	var activityReportGroups, e = me.buildReport(repositoryLogEntries)
	common.AssertError(e)

	if false {
		var pluginFilePath = common.ExecutableFilePath + "/" + me.Configuration.Plugin
		var pluginRunner = PluginRunner{FilePath: pluginFilePath}
		var storyEntries, pluginError = pluginRunner.Run(git_stories_api.StoriesRequest{
			LogEntries: repositoryLogEntries,
			TimeZone:   me.Request.TimeZone,
		})
		if nil != pluginError {
			me.Update(&WebTask{Error: pluginError.Error()})
			return
		}
		common.Use(storyEntries)
	}

	var task WebTask
	task.ActivityReportGroups = activityReportGroups
	if nil == task.ActivityReportGroups {
		// Avoid nil value because nil means that the task is not finished yet
		task.ActivityReportGroups = make(ActivityReportGroups, 0)
	}
	if me.Request.LengthLimit > 0 && len(task.ActivityReportGroups) > me.Request.LengthLimit {
		task.ActivityReportGroups = task.ActivityReportGroups[0:me.Request.LengthLimit]
	}
	me.Update(&task)
}

func (me *ReportGenerator) buildReport(repositoryLogEntries git_stories_api.RepositoryLogEntries) (
	/*Returns*/ groups ActivityReportGroups, e error,
) {
	var timeZone, locationError = time.LoadLocation(me.Request.TimeZone)
	if locationError != nil {
		return nil, locationError
	}
	var reportByDate map[int]*ActivityReportGroup = make(map[int]*ActivityReportGroup)
	for _, entry := range repositoryLogEntries {
		var entryTime = entry.Time.In(timeZone)
		var isMerge = len(entry.Parents) > 1
		var report = me.getOrCreate(reportByDate, entryTime)

		var activity ActivityReport
		if isMerge {
			// I assume that normally people do only a little work for a merge commit
			// Therefore merges are not included into the usual activity
			activity.Points = 1
		} else {
			activity.ReadRepositoryLogEntry(entry)
		}
		report.Aggregate(entry.AuthorName, &activity)
	}

	for _, reportGroup := range reportByDate {
		groups = append(groups, reportGroup)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Time.Before(groups[j].Time)
	})
	return
}

func (me *ReportGenerator) getOrCreate(reportByDate map[int]*ActivityReportGroup, entryTime time.Time) (report *ActivityReportGroup) {
	var year, month, day = entryTime.Date()
	var dateHash = GetDateHash(year, month, day)
	report = reportByDate[dateHash]
	if nil == report {
		report = &ActivityReportGroup{
			Period: time.Hour * 24,
			Time:   time.Date(year, month, day, 0, 0, 0, 0, entryTime.Location()),
		}
		reportByDate[dateHash] = report
	}
	return
}
