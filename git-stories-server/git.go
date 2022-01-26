package main

import (
	"errors"
	"strconv"
	"strings"

	git_stories_api "github.com/hinst/git-stories-api"
)

const GitRootNodeHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

func GitParseSummaryCount(text string) int {
	if text == "-" {
		return git_stories_api.DiffSummaryBinary
	} else {
		var result, e = strconv.Atoi(text)
		if e != nil {
			return git_stories_api.DiffSummaryError
		}
		return result
	}
}

func IsSpaceOrTab(r rune) bool {
	return r == ' ' || r == '\t'
}

func IsNotSpaceOrTab(r rune) bool {
	return !IsSpaceOrTab(r)
}

func GitSplitDiffSummaryLine(line string) (parts []string, e error) {
	var lookForInsertionCountStart bool = true
	var insertionCountRunes []rune
	var lookForInsertionCountEnd bool = false
	var lookForDeletionCountStart bool = false
	var deletionCountRunes []rune
	var lookForDeletionCountEnd bool = false
	var lookForFilePathStart bool = false
	var filePathBuilder strings.Builder
	var lookForFilePathEnd bool = false
	for _, r := range line {
		if lookForInsertionCountStart {
			if IsNotSpaceOrTab(r) {
				lookForInsertionCountStart = false
				lookForInsertionCountEnd = true
				insertionCountRunes = append(insertionCountRunes, r)
			}
		} else if lookForInsertionCountEnd {
			if IsSpaceOrTab(r) {
				lookForInsertionCountEnd = false
				lookForDeletionCountStart = true
			} else {
				insertionCountRunes = append(insertionCountRunes, r)
			}
		} else if lookForDeletionCountStart {
			if IsNotSpaceOrTab(r) {
				lookForDeletionCountStart = false
				lookForDeletionCountEnd = true
				deletionCountRunes = append(deletionCountRunes, r)
			}
		} else if lookForDeletionCountEnd {
			if IsSpaceOrTab(r) {
				lookForDeletionCountEnd = false
				lookForFilePathStart = true
			} else {
				deletionCountRunes = append(deletionCountRunes, r)
			}
		} else if lookForFilePathStart {
			if IsNotSpaceOrTab(r) {
				lookForFilePathStart = false
				lookForFilePathEnd = true
				filePathBuilder.WriteRune(r)
			}
		} else if lookForFilePathEnd {
			filePathBuilder.WriteRune(r)
		}
	}
	if nil == insertionCountRunes {
		return nil, errors.New("unable to find insertion count")
	}
	if nil == deletionCountRunes {
		return nil, errors.New("unable to find deletion count")
	}
	return []string{
			string(insertionCountRunes),
			string(deletionCountRunes),
			filePathBuilder.String(),
		},
		nil
}
