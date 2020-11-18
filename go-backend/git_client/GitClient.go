package git_client

import (
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

func (this *GitClient) GetCommandName() string {
	if this.commandName != "" {
		return this.commandName
	} else {
		return "git"
	}
}
