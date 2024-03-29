package git_secretary

import (
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	git_stories_api "github.com/hinst/git-stories-api"
	common "github.com/hinst/go-common"
)

type GitClient struct {
	directory       string
	commandName     string
	debugLogEnabled bool
}

func CreateGitClient(directory string) *GitClient {
	return &GitClient{directory: directory}
}

func (gitClient *GitClient) Run(args []string) (string, error) {
	var commandName = gitClient.GetCommandName()
	if gitClient.debugLogEnabled {
		log.Println("GitClient.Run: " + commandName + " " + strings.Join(args, " "))
	}
	var command = exec.Command(commandName, args...)
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

func (gitClient *GitClient) ReadLog(lengthLimit int) ([]RepositoryLogEntryHeader, error) {
	var args = []string{"log", "--all", "--format=\"%H %P\""}
	if lengthLimit > 0 {
		args = append(args, "-n", strconv.Itoa(lengthLimit))
	} else {
		args = append(args, "-n", strconv.Itoa(math.MaxInt))
	}
	args = append(args, "HEAD")
	var outputText, e = gitClient.Run(args)
	if e == nil {
		var lines = strings.Split(outputText, "\n")
		var rows []RepositoryLogEntryHeader
		for i := 0; i < len(lines); i++ {
			var line = lines[i]
			line = strings.Trim(line, "\"")
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				var logEntryRow RepositoryLogEntryHeader
				logEntryRow.ParseGitLine(line)
				rows = append(rows, logEntryRow)
			}
		}
		return rows, nil
	} else {
		return nil, e
	}
}

var ReadCommitInfoItems = []string{
	GitShowFormat.AuthorTime,
	GitShowFormat.AuthorNameMapped,
}
var ReadCommitInfoFormat = strings.Join(ReadCommitInfoItems, GitShowFormat.NewLine)

func (gitClient *GitClient) ReadCommitInfo(row *git_stories_api.RepositoryLogEntry) error {
	var outputText, runError = gitClient.Run([]string{
		"show",
		"--format=" + ReadCommitInfoFormat,
		"--quiet",
		row.CommitHash,
	})
	if runError != nil {
		return runError
	}
	var outputParts = strings.Split(outputText, "\n")
	var authorTimeText = outputParts[0]
	var authorNameText = outputParts[1]
	var second, numberError = strconv.ParseInt(strings.TrimSpace(authorTimeText), 10, 64)
	if numberError != nil {
		return common.CreateException("ReadCommitDate", numberError)
	}
	row.Time = time.Unix(second, 0)
	row.AuthorName = strings.TrimSpace(authorNameText)
	return nil
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
			var parts, e = GitSplitDiffSummaryLine(line)
			if nil != e {
				continue
			}
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

func (gitClient *GitClient) ReadDetailedLog(lengthLimit int) ([]git_stories_api.RepositoryLogEntry, error) {
	var logEntries, readError = gitClient.ReadLog(lengthLimit)
	if nil != readError {
		return nil, readError
	}
	var rows []git_stories_api.RepositoryLogEntry
	for _, entry := range logEntries {
		var row, e = gitClient.ReadDetailedLogEntryRow(entry)
		if nil != e {
			return nil, e
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (gitClient *GitClient) ReadDetailedLogEntryRow(header RepositoryLogEntryHeader) (row git_stories_api.RepositoryLogEntry, e error) {
	var parentHashes []string
	if len(header.ParentHashes) > 0 {
		parentHashes = header.ParentHashes
	} else {
		parentHashes = []string{GitRootNodeHash}
	}
	var parentInfos []git_stories_api.ParentInfoEntry
	for _, parentHash := range parentHashes {
		var diffSummary, diffSummaryError = gitClient.ReadDiffSummary(parentHash, header.CommitHash)
		if nil != diffSummaryError {
			return row, diffSummaryError
		}
		var parentInfo = git_stories_api.ParentInfoEntry{
			CommitHash: parentHash,
			DiffRows:   diffSummary,
		}
		parentInfos = append(parentInfos, parentInfo)
	}
	row = git_stories_api.RepositoryLogEntry{
		CommitHash: header.CommitHash,
		Parents:    parentInfos,
	}
	var readCommitInfoError = gitClient.ReadCommitInfo(&row)
	if nil != readCommitInfoError {
		return row, readCommitInfoError
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

func (gitClient *GitClient) SetDebugLogEnabled(enabled bool) *GitClient {
	gitClient.debugLogEnabled = enabled
	return gitClient
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
