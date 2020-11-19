package git_client

import (
	"strconv"
	"time"
	"os/exec"
	"strings"
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
	var outputText, e = this.Run([]string{"log", "--format=%H %P"})
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
	if (numberError != nil) {
		return time.Unix(0, 0), numberError
	}
	return time.Unix(second, 0), nil
}

func (this *GitClient) GetCommandName() string {
	if this.commandName != "" {
		return this.commandName
	} else {
		return "git"
	}
}
