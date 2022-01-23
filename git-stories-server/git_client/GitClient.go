package git_client

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
)

type GitClient struct {
	directory   string
	commandName string
}

func CreateGitClient(directory string) *GitClient {
	return &GitClient{directory: directory}
}

func (gitClient *GitClient) Run(args []string) (string, error) {
	var command = exec.Command(gitClient.GetCommandName(), args...)
	command.Dir = gitClient.directory
	var output, e = command.CombinedOutput()
	var outputText = string(output)
	if e == nil {
		return outputText, e
	} else {
		var errorMessage = "Error running git " + strings.Join(args, " ") + "\n" + outputText
		return "", common.CreateException(errorMessage, e)
	}
}

func (gitClient *GitClient) ReadLog(lengthLimit int) ([]LogEntryRow, error) {
	var args = []string{"log", "--format=\"%H %P\""}
	if lengthLimit > 0 {
		args = append(args, "-n", strconv.Itoa(lengthLimit))
	}
	args = append(args, "HEAD")
	var outputText, e = gitClient.Run(args)
	if e == nil {
		var lines = strings.Split(outputText, "\n")
		var rows []LogEntryRow
		for i := 0; i < len(lines); i++ {
			var line = lines[i]
			line = strings.Trim(line, "\"")
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				var logEntryRow LogEntryRow
				logEntryRow.ParseGitLine(line)
				rows = append(rows, logEntryRow)
			}
		}
		return rows, nil
	} else {
		return nil, e
	}
}

func (gitClient *GitClient) ReadCommitDate(commitHash string) (time.Time, error) {
	var outputText, runError = gitClient.Run([]string{"show", "--format=%at", "--quiet", commitHash})
	if runError != nil {
		return time.Unix(0, 0), runError
	}
	var second, numberError = strconv.ParseInt(strings.TrimSpace(outputText), 10, 64)
	if numberError != nil {
		return time.Unix(0, 0), common.CreateException("ReadCommitDate", numberError)
	}
	return time.Unix(second, 0), nil
}

func checkSpaceOrTab(r rune) bool {
	return r == ' ' || r == '\t'
}

func (gitClient *GitClient) ReadDiffSummary(commitHash1, commitHash2 string) ([]git_stories_api.DiffSummaryRow, error) {
	var outputText, runError = gitClient.Run([]string{"diff", "--numstat", commitHash1, commitHash2})
	if runError != nil {
		return nil, runError
	}
	var lines = strings.Split(outputText, "\n")
	var rows []git_stories_api.DiffSummaryRow
	for i := 0; i < len(lines); i++ {
		var line = strings.TrimSpace(lines[i])
		if len(line) > 0 {
			var parts = strings.FieldsFunc(line, checkSpaceOrTab)
			var row = git_stories_api.DiffSummaryRow{
				InsertionCount: diffSummaryPartToLong(parts[0]),
				DeletionCount:  diffSummaryPartToLong(parts[1]),
				FilePath:       parts[2],
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (gitClient *GitClient) ReadDetailedLog(lengthLimit int) ([]git_stories_api.DetailedLogEntryRow, error) {
	var logEntries, readError = gitClient.ReadLog(lengthLimit)
	if nil != readError {
		return nil, readError
	}
	var rows []git_stories_api.DetailedLogEntryRow
	for _, entry := range logEntries {
		var row, e = gitClient.ReadDetailedLogEntryRow(entry)
		if nil != e {
			return nil, e
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (gitClient *GitClient) ReadDetailedLogEntryRow(logEntry LogEntryRow) (git_stories_api.DetailedLogEntryRow, error) {
	var commitDate, commitDateError = gitClient.ReadCommitDate(logEntry.CommitHash)
	if nil != commitDateError {
		return git_stories_api.DetailedLogEntryRow{}, commitDateError
	}
	var parentHashes []string
	if len(logEntry.ParentHashes) > 0 {
		parentHashes = logEntry.ParentHashes
	} else {
		parentHashes = []string{GitRootNodeHash}
	}
	var parentInfos []git_stories_api.ParentInfoEntry
	for _, parentHash := range parentHashes {
		var diffSummary, diffSummaryError = gitClient.ReadDiffSummary(parentHash, logEntry.CommitHash)
		if nil != diffSummaryError {
			return git_stories_api.DetailedLogEntryRow{}, diffSummaryError
		}
		var parentInfo = git_stories_api.ParentInfoEntry{
			CommitHash:  parentHash,
			DiffSummary: diffSummary,
		}
		parentInfos = append(parentInfos, parentInfo)
	}
	var row = git_stories_api.DetailedLogEntryRow{
		CommitHash: logEntry.CommitHash,
		Time:       commitDate,
		Parents:    parentInfos,
	}
	return row, nil
}

func (gitClient *GitClient) GetCommandName() string {
	if gitClient.commandName != "" {
		return gitClient.commandName
	} else {
		return "git"
	}
}

func diffSummaryPartToLong(text string) int {
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
