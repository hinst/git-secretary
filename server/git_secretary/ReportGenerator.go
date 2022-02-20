package git_secretary

import (
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
	me.buildReport(repositoryLogEntries)

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

	var task WebTask
	task.StoryEntries = storyEntries
	if nil == task.StoryEntries {
		// Avoid nil value because nil means that the task is not finished yet
		task.StoryEntries = make([]git_stories_api.StoryEntryChangeset, 0)
	}
	if me.Request.LengthLimit > 0 && len(task.StoryEntries) > me.Request.LengthLimit {
		task.StoryEntries = task.StoryEntries[0:me.Request.LengthLimit]
	}
	me.Update(&task)
}

func (me *ReportGenerator) buildReport(repositoryLogEntries git_stories_api.RepositoryLogEntries) (
	/*Returns*/ reportEntries ReportEntries, e error,
) {
	var timeZone, locationError = time.LoadLocation(me.Request.TimeZone)
	if locationError != nil {
		return nil, locationError
	}
	var reportByDate map[int]*ReportEntry = make(map[int]*ReportEntry)
	for _, entry := range repositoryLogEntries {
		var entryTime = entry.Time.In(timeZone)
		var year, month, day = entryTime.Date()
		var date = GetDateHash(year, month, day)
		var report = reportByDate[date]
		if nil == report {
			report = &ReportEntry{Time: time.Date(year, month, day, 0, 0, 0, 0, timeZone)}
		}
		var isMerge = len(entry.Parents) > 1
		var entryReport = &ReportEntry{
			Time:   entryTime,
			Period: time.Hour * 24,
		}
		entryReport.Activity.Points = 1
		// I assume that normally people do only a little work for a merge commit
		// Therefore the user gets only 1 activity point for a merge
		if !isMerge {
			me.addActivity(&entryReport.Activity, entry)
		}
	}
	return
}

func (me *ReportGenerator) addActivity(activity *ReportActivityEntry, source *git_stories_api.RepositoryLogEntry) {
	activity.ChangesetCount += 1
	for _, parent := range source.Parents {
		activity.ChangedFileCount += len(parent.DiffRows)
		for _, diff := range parent.DiffRows {
			activity.ChangedFileCount += 1
			activity.InsertionCount += MaxOfTwoInts(0, diff.InsertionCount)
			activity.DeletionCount += MaxOfTwoInts(0, diff.DeletionCount)
		}
	}
}
