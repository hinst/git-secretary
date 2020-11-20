package git_client

import (
	"git-stories/common"
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
	var command = exec.Command(this.GetCommandName(), args...)
	command.Dir = this.directory
	var output, e = command.CombinedOutput()
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
	var second, numberError = strconv.ParseInt(strings.TrimSpace(outputText), 10, 64)
	if numberError != nil {
		return time.Unix(0, 0), common.WrapError("ReadCommitDate", numberError)
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
		var row, e = this.ReadDetailedLogEntryRow(entry)
		if nil != e {
			return nil, e
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (this *GitClient) ReadDetailedLogEntryRow(logEntry LogEntryRow) (DetailedLogEntryRow, error) {
	var commitDate, commitDateError = this.ReadCommitDate(logEntry.CommitHash)
	if nil != commitDateError {
		return DetailedLogEntryRow{}, commitDateError
	}
	var parentHashes []string
	if len(logEntry.ParentHashes) > 0 {
		parentHashes = logEntry.ParentHashes
	} else {
		parentHashes = []string{GitRootNodeHash}
	}
	var parentInfos []ParentInfoEntry
	for _, parentHash := range parentHashes {
		var diffSummary, diffSummaryError = this.ReadDiffSummary(parentHash, logEntry.CommitHash)
		if nil != diffSummaryError {
			return DetailedLogEntryRow{}, diffSummaryError
		}
		var parentInfo = ParentInfoEntry{
			CommitHash:  parentHash,
			DiffSummary: diffSummary,
		}
		parentInfos = append(parentInfos, parentInfo)
	}
	var row = DetailedLogEntryRow{
		CommitHash: logEntry.CommitHash,
		Time:       commitDate,
		Parents:    parentInfos,
	}
	return row, nil
}

func (this *GitClient) GetCommandName() string {
	if this.commandName != "" {
		return this.commandName
	} else {
		return "git"
	}
}
