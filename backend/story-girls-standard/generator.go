package main

import (
	"bytes"
	"github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	"hash/fnv"
	"html/template"
	"strings"
)

type StoryEntry git_stories_api.StoryEntry

type Generator struct {
	Entries []git_stories_api.DetailedLogEntryRow
}

var girlNames = loadGirlNames()

func loadGirlNames() []string {
	var girlNamesData, assetError = Asset("resources/girl_names.txt")
	common.AssertError(assetError)
	var names = strings.Split(string(girlNamesData), "\n")
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
	var allStoryEntries []StoryEntry
	for _, logEntry := range this.Entries {
		var storyEntries = this.generateEntries(logEntry)
		allStoryEntries = append(allStoryEntries, storyEntries...)
	}
	return allStoryEntries
}

func (this *Generator) generateEntries(entry git_stories_api.DetailedLogEntryRow) []StoryEntry {
	var entries []StoryEntry
	for _, parent := range entry.Parents {
		for _, diffSummaryRow := range parent.DiffSummary {
			var description = this.generateStoryDescriptionFromDiffRow(diffSummaryRow)
			if "" != description {
				var storyEntry = StoryEntry{
					Time:        entry.Time,
					CommitHash:  entry.CommitHash,
					ParentHash:  parent.CommitHash,
					Description: description,
				}
				entries = append(entries, storyEntry)
			}
		}
	}
	return entries
}

func (this *Generator) generateStoryDescriptionFromDiffRow(row git_stories_api.DiffSummaryRow) string {
	var nameHash = getHashFromString(row.FilePath)
	var nameIndex = nameHash % uint32(len(girlNames))
	var characterName = girlNames[nameIndex]
	var actionHash = getHashFromString(row.FilePath)
	var actionTemplate string
	actionTemplate = this.getActionTemplate(row, actionHash, actionTemplate)
	var actionArgs = getActionArgs(row)
	actionArgs.Name = characterName
	var parsedTemplate, templateError = template.New("").Parse(actionTemplate)
	common.AssertError(templateError)
	var description = getStringFromTemplate(parsedTemplate, actionArgs)
	return description
}

func (this *Generator) getActionTemplate(row git_stories_api.DiffSummaryRow, actionHash uint32, actionTemplate string) string {
	if row.InsertionCount > 0 && row.DeletionCount == 0 {
		var actionIndex = actionHash % uint32(len(Actions.InsertionActions))
		actionTemplate = Actions.InsertionActions[actionIndex]
	} else if row.InsertionCount == 0 && row.DeletionCount > 0 {
		var actionIndex = actionHash % uint32(len(Actions.DeletionActions))
		actionTemplate = Actions.DeletionActions[actionIndex]
	} else if row.InsertionCount > 0 && row.DeletionCount > 0 {
		var actionIndex = actionHash % uint32(len(Actions.CombinedActions))
		actionTemplate = Actions.CombinedActions[actionIndex]
	} else {
		actionTemplate = ""
	}
	return actionTemplate
}

func getHashFromString(s string) uint32 {
	var hasher = fnv.New32a()
	hasher.Write([]byte(s))
	return hasher.Sum32()
}

func getActionArgs(row git_stories_api.DiffSummaryRow) ActionArgs {
	var actionArgs ActionArgs
	actionArgs.InsertionCount = row.InsertionCount
	if row.InsertionCount > 1 {
		actionArgs.SI = "s"
	}
	actionArgs.DeletionCount = row.DeletionCount
	if row.DeletionCount > 1 {
		actionArgs.SD = "s"
	}
	return actionArgs
}

func getStringFromTemplate(t *template.Template, args interface{}) string {
	var buffer bytes.Buffer
	var templateError = t.Execute(&buffer, args)
	common.AssertError(templateError)
	return buffer.String()
}
