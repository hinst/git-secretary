package main

import (
	"bufio"
	"encoding/json"
	"github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	"io/ioutil"
	"os"
)

func main() {
	var reader = bufio.NewReaderSize(os.Stdin, 512)
	var input, inputError = ioutil.ReadAll(reader)
	common.AssertError(inputError)
	var entries []git_stories_api.DetailedLogEntryRow
	var jsonError = json.Unmarshal(input, &entries)
	common.AssertError(jsonError)
	var generator = Generator{Entries: entries}
	var stories = generator.Generate()
	var storiesJson, jsonWriterError = json.Marshal(stories)
	common.AssertError(jsonWriterError)
	var _, writeError = os.Stdout.Write(storiesJson)
	common.AssertError(writeError)
}
