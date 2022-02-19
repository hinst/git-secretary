package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/hinst/go-common"
)

type FilePicker struct {
	WebPath string
}

type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
}

func (picker *FilePicker) Initialize(wrapper func(pattern string, handler http.HandlerFunc)) {
	wrapper(picker.WebPath+"/fileList", picker.GetFileList)
}

func (picker *FilePicker) GetFileList(responseWriter http.ResponseWriter, request *http.Request) {
	var fileInfos []FileInfo
	var directory = request.URL.Query().Get("directory")
	if len(directory) == 0 {
		if CheckWindows() {
			var driveList = picker.getDriveList()
			for _, driveLetter := range driveList {
				fileInfos = append(fileInfos, FileInfo{
					Path: driveLetter, Name: driveLetter, IsDirectory: true,
				})
			}
		} else {
			directory = "/"
		}
	}
	if len(directory) > 0 {
		var files, fileError = os.ReadDir(picker.getDirectory(directory))
		if fileError != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			responseWriter.Write([]byte("Unable to read directory \"" + directory + "\""))
			return
		}
		for _, fileInfo := range files {
			fileInfos = append(fileInfos, FileInfo{
				Path:        directory + string(os.PathSeparator) + fileInfo.Name(),
				Name:        fileInfo.Name(),
				IsDirectory: fileInfo.IsDir(),
			})
		}
	}
	var data, e = json.Marshal(fileInfos)
	common.AssertError(common.CreateExceptionIf("error: write JSON", e))
	responseWriter.Write(data)
}

func (picker *FilePicker) getDriveList() (drives []string) {
	var command = exec.Command("wmic", "logicaldisk", "get", "deviceid")
	var output, e = command.CombinedOutput()
	common.AssertError(common.CreateExceptionIf("Unable to read list of logical drives", e))
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

func (picker *FilePicker) getDirectory(directory string) string {
	// Workaround for directory behavior on Windows
	var isDriveLetter = len(directory) == 2 && directory[1] == ':'
	if isDriveLetter {
		return directory + string(os.PathSeparator)
	} else {
		return directory
	}
}
