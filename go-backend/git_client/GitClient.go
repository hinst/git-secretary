package git_client

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type GitClient struct {
	directory   string
	commandName string
}

func CreateGitClient(directory string) *GitClient {
	return &GitClient{directory: directory}
}

func (this *GitClient) Run(args []string) (string, error) {
	var output, e = exec.Command(this.GetCommandName(), args...).CombinedOutput()
	if e == nil {
		var outputText = string(output)
		return outputText, e
	} else {
		return "", e
	}
}

func (this *GitClient) ReadLog() ([]LogEntryRow, error) {
	var outputText, e = this.Run([]string{"log", "HEAD", "--format=%H %P"})
	if e == nil {
		var lines = strings.Split(outputText, "\n")
		var rows []LogEntryRow
		for i := 0; i < len(lines); i++ {
			var line = strings.TrimSpace(lines[i])
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

func (this *GitClient) ReadCommitDate(commitHash string) (time.Time, error) {
	var outputText, runError = this.Run([]string{"show", "--format=%at", "--quiet", commitHash})
	if runError != nil {
		return time.Unix(0, 0), runError
	}
	var text = string(outputText)
	var second, numberError = strconv.ParseInt(text, 10, 64)
	if numberError != nil {
		return time.Unix(0, 0), numberError
	}
	return time.Unix(second, 0), nil
}

func checkSpaceOrTab(r rune) bool {
	return r == ' ' || r == '\t'
}

func (this *GitClient) ReadDiffSummary(commitHash1, commitHash2 string) ([]DiffSummaryRow, error) {
	var outputText, runError = this.Run([]string{"diff", "--numstat", commitHash1, commitHash2})
	if runError != nil {
		return nil, runError
	}
	var lines = strings.Split(outputText, "\n")
	var rows []DiffSummaryRow
	for i := 0; i < len(lines); i++ {
		var line = strings.TrimSpace(lines[i])
		if len(line) > 0 {
			var parts = strings.FieldsFunc(line, checkSpaceOrTab)
			var row = DiffSummaryRow{
				InsertionCount: diffSummaryPartToLong(parts[0]),
				DeletionCount:  diffSummaryPartToLong(parts[1]),
				FilePath:       parts[2],
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (this *GitClient) ReadAllDetailedLog() ([]DetailedLogEntryRow, error) {
	var logEntries, readError = this.ReadLog()
	if nil != readError {
		return nil, readError
	}
	var rows []DetailedLogEntryRow
	for _, entry := range logEntries {
		var commitDate, commitDateError = this.ReadCommitDate(entry.CommitHash)
		if nil != commitDateError {
			return nil, commitDateError
		}
		var diffSummary, diffSummaryError = this.ReadDiffSummary(entry.ParentHashes[0], entry.CommitHash)
		if nil != diffSummaryError {
			return nil, diffSummaryError
		}
		var row = DetailedLogEntryRow{
			LogEntry:    entry,
			Time:        commitDate,
			DiffSummary: diffSummary,
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (this *GitClient) GetCommandName() string {
	if this.commandName != "" {
		return this.commandName
	} else {
		return "git"
	}
}
