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

func (runner *PluginRunner) Run(rows git_stories_api.RepositoryLogEntries) ([]git_stories_api.StoryEntryChangeset, error) {
	var rowsBytes, jsonWriteError = json.Marshal(rows)
	if nil != jsonWriteError {
		return nil, common.CreateException("Cannot to write json", jsonWriteError)
	}
	var pluginFilePath = runner.PluginFilePath
	if CheckWindows() && !strings.HasSuffix(pluginFilePath, WINDOWS_EXECUTABLE_FILE_EXTENSION) {
		pluginFilePath += WINDOWS_EXECUTABLE_FILE_EXTENSION
	}
	var command = exec.Command(pluginFilePath)
	var writer, writerError = command.StdinPipe()
	if nil != writerError {
		return nil, common.CreateException("Cannot write input for the plug-in", writerError)
	}
	var output, pluginOutputError = command.StdoutPipe()
	if nil != pluginOutputError {
		return nil, common.CreateException("Cannot open output of the plug-in", pluginOutputError)
	}
	var errorOutput, errorOutputError = command.StderrPipe()
	if nil != errorOutputError {
		return nil, common.CreateException("Cannot open stderr of the plug-in", errorOutputError)
	}
	var commandStartError = command.Start()
	if nil != commandStartError {
		return nil, common.CreateException("Cannot start command", commandStartError)
	}
	var bufferedWriter = bufio.NewWriter(writer)
	var _, bufferedWriteError = bufferedWriter.Write(rowsBytes)
	if nil != bufferedWriteError {
		return nil, common.CreateException("Cannot write buffered", bufferedWriteError)
	}
	var flushError = bufferedWriter.Flush()
	if nil != flushError {
		return nil, common.CreateException("Cannot flush writer", flushError)
	}
	var closeError = writer.Close()
	if nil != closeError {
		return nil, common.CreateException("Cannot close writer", closeError)
	}
	var outputData, outputError = ioutil.ReadAll(bufio.NewReader(output))
	if nil != outputError {
		return nil, common.CreateException("Cannot read output of the plug-in", outputError)
	}
	var errorData, errorError = ioutil.ReadAll(bufio.NewReader(errorOutput))
	if nil != errorError {
		return nil, common.CreateException("Cannot read error output of the plug-in", errorError)
	}
	println(string(errorData))
	var storyEntries []git_stories_api.StoryEntryChangeset
	var jsonError = json.Unmarshal(outputData, &storyEntries)
	if nil != jsonError {
		return nil, common.CreateException("Cannot decode output of the plug-in as JSON", jsonError)
	}
	return storyEntries, nil
}
