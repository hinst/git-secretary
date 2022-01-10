package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hinst/go-common"
)

type FilePicker struct {
	WebPath string
}

type FileInfo struct {
	name        string
	isDirectory bool
}

func (picker *FilePicker) Initialize(wrapper func(pattern string, handler http.HandlerFunc)) {
	wrapper(picker.WebPath+"/fileList", picker.GetFileList)
}

func (picker *FilePicker) GetFileList(responseWriter http.ResponseWriter, request *http.Request) {
	var fileInfos []FileInfo
	var directory = request.URL.Query().Get("directory")
	if len(directory) == 0 {
		if runtime.GOOS == "windows" {
			var driveList = picker.getDriveList()
			for _, driveLetter := range driveList {
				fileInfos = append(fileInfos, FileInfo{name: driveLetter + "\\", isDirectory: true})
			}
		} else {
			directory = "/"
		}
	}
	if len(directory) > 0 {
		filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
			if info != nil {
				fileInfos = append(fileInfos, FileInfo{name: path, isDirectory: info.IsDir()})
			}
			return nil
		})
	}
	var data, error = json.Marshal(fileInfos)
	common.AssertWrapped(error, "error: write JSON")
	responseWriter.Write(data)
}

func (picker *FilePicker) getDriveList() (drives []string) {
	var command = exec.Command("wmic", "logicaldisk", "get", "deviceid")
	var output, e = command.CombinedOutput()
	common.AssertWrapped(e, "Unable to read list of logical drives")
	var outputText = string(output)
	var lines = strings.Split(outputText, "\n")
	for i, line := range lines {
		if i > 0 {
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				drives = append(drives, line)
			}
		}
	}
	return
}
