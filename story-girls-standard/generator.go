package main

import (
	"bytes"
	_ "embed"
	"hash/fnv"
	"html/template"
	"strings"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
)

type StoryEntry git_stories_api.StoryEntry

type Generator struct {
	Entries []git_stories_api.DetailedLogEntryRow
}

//go:embed resources/girl_names.txt
var girlNamesData string
var girlNames = loadGirlNames()

func loadGirlNames() []string {
	var names = strings.Split(girlNamesData, "\n")
	var processedNames []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if len(name) > 0 {
			processedNames = append(processedNames, name)
		}
	}
	return processedNames
}

func (me *Generator) Generate() []StoryEntry {
	var allStoryEntries []StoryEntry
	for _, logEntry := range me.Entries {
		var storyEntries = me.generateEntries(logEntry)
		allStoryEntries = append(allStoryEntries, storyEntries...)
	}
	return allStoryEntries
}

func (me *Generator) generateEntries(entry git_stories_api.DetailedLogEntryRow) []StoryEntry {
	var entries []StoryEntry
	for _, parent := range entry.Parents {
		for _, diffSummaryRow := range parent.DiffSummary {
			var description = me.generateStoryDescriptionFromDiffRow(entry.CommitHash, diffSummaryRow)
			if len(description) > 0 {
				var storyEntry = StoryEntry{
					Time:           entry.Time,
					CommitHash:     entry.CommitHash,
					ParentHash:     parent.CommitHash,
					Description:    description,
					SourceFilePath: diffSummaryRow.FilePath,
				}
				entries = append(entries, storyEntry)
			}
		}
	}
	return entries
}

func (me *Generator) generateStoryDescriptionFromDiffRow(commitHash string, row git_stories_api.DiffSummaryRow) string {
	var nameHash = getHashFromString(row.FilePath)
	var nameIndex = nameHash % uint32(len(girlNames))
	var characterName = girlNames[nameIndex]
	var actionHash = getHashFromString(commitHash)
	var actionTemplate = me.getActionTemplate(row, actionHash)
	var actionArgs = getActionArgs(row)
	actionArgs.Name = characterName
	var parsedTemplate, templateError = template.New("").Parse(actionTemplate)
	common.AssertError(templateError)
	var description = getStringFromTemplate(parsedTemplate, actionArgs)
	return description
}

func (me *Generator) getActionTemplate(row git_stories_api.DiffSummaryRow, actionHash uint32) (actionTemplate string) {
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
	var _, e = hasher.Write([]byte(s))
	common.AssertError(e)
	return hasher.Sum32()
}

func getActionArgs(row git_stories_api.DiffSummaryRow) ActionArgs {
	var actionArgs ActionArgs
	actionArgs.InsertionCount = row.InsertionCount
	if row.InsertionCount > 1 {
		actionArgs.IS = "s"
		actionArgs.IES = "es"
	}
	actionArgs.DeletionCount = row.DeletionCount
	if row.DeletionCount > 1 {
		actionArgs.DS = "s"
		actionArgs.DES = "es"
	}
	actionArgs.HeShe = "she"
	return actionArgs
}

func getStringFromTemplate(t *template.Template, args interface{}) string {
	var buffer bytes.Buffer
	var templateError = t.Execute(&buffer, args)
	common.AssertError(templateError)
	return buffer.String()
}
