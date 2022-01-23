package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
)

type PluginRunner struct {
	PluginFilePath string
}

func (runner *PluginRunner) Run(rows git_stories_api.DetailedLogEntryRows) []git_stories_api.StoryEntry {
	var rowsBytes, jsonWriteError = json.Marshal(rows)
	common.AssertError(jsonWriteError)
	var pluginFilePath = runner.PluginFilePath
	if CheckWindows() && !strings.HasSuffix(pluginFilePath, WINDOWS_EXECUTABLE_FILE_EXTENSION) {
		pluginFilePath += WINDOWS_EXECUTABLE_FILE_EXTENSION
	}
	var command = exec.Command(pluginFilePath)
	var writer, writerError = command.StdinPipe()
	common.AssertError(writerError)
	var output, pluginOutputError = command.StdoutPipe()
	common.AssertError(pluginOutputError)
	common.AssertError(command.Start())
	var bufferedWriter = bufio.NewWriter(writer)
	var _, bufferedWriteError = bufferedWriter.Write(rowsBytes)
	common.AssertError(bufferedWriteError, bufferedWriter.Flush())
	var closeError = writer.Close()
	common.AssertError(closeError)
	var outputData, outputError = ioutil.ReadAll(bufio.NewReader(output))
	common.AssertError(outputError)
	var storyEntries []git_stories_api.StoryEntry
	json.Unmarshal(outputData, &storyEntries)
	return storyEntries
}
