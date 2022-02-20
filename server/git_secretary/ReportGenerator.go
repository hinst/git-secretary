package git_secretary

import (
	"os"

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
		me.Update(&WebTask{Total: total, Done: done})
	})
	var rows, gitError = gitClient.ReadDetailedLog(0)
	if nil != gitError {
		me.Update(&WebTask{Error: gitError.Error()})
		return
	}
	var workingDirectory, getwdError = os.Getwd()
	common.AssertError(getwdError)
	var pluginFilePath = workingDirectory + "/" + me.Configuration.Plugin
	var pluginRunner = PluginRunner{PluginFilePath: pluginFilePath}
	var storyEntries, pluginError = pluginRunner.Run(git_stories_api.StoriesRequest{
		LogEntries: rows,
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
