package main

import (
	"github.com/hinst/git-stories-api"
	"strings"
)

type StoryEntry git_stories_api.StoryEntry

type Generator struct {
	Entries []git_stories_api.DetailedLogEntryRow
}

var girlNames = loadGirlNames()

func loadGirlNames() []string {
	var names = strings.Split(string(_resourcesGirl_namesTxt), "\n")
	var processedNames []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if len(name) > 0 {
			processedNames = append(processedNames, name)
		}
	}
	return processedNames
}

func (this *Generator) Generate() []StoryEntry {
	var storyEntries []StoryEntry
	for _, logEntry := range this.Entries {
		var storyEntries = this.generateEntries(logEntry)
		storyEntries = append(storyEntries, storyEntries...)
	}
	return storyEntries
}

func (this *Generator) generateEntries(entry git_stories_api.DetailedLogEntryRow) []StoryEntry {
}
