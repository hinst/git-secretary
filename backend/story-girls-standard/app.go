package main

import (
	"bufio"
	"encoding/json"
	"github.com/hinst/git-stories-api"
	"io/ioutil"
	"os"
)

func main() {
	var reader = bufio.NewReader(os.Stdin)
	var input, inputError = ioutil.ReadAll(reader)
	if nil != inputError {
		panic(inputError)
	}
	var entries []git_stories_api.DetailedLogEntryRow
	var jsonError = json.Unmarshal(input, &entries)
	if nil != jsonError {
		panic(jsonError)
	}
}
