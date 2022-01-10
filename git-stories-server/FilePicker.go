package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/hinst/go-common"
)

type FilePicker struct {
}

type FileInfo struct {
	name        string
	isDirectory bool
}

func (picker *FilePicker) GetFileList(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query().Get("directory")
	var fileInfos []FileInfo
	filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		fileInfos = append(fileInfos, FileInfo{name: path, isDirectory: info.IsDir()})
		return nil
	})
	var data, error = json.Marshal(fileInfos)
	common.AssertWrapped(error, "error: write JSON")
	responseWriter.Write(data)
}
