package main

import (
	"encoding/json"
	"git-stories-server/git_client"
	"github.com/hinst/go-common"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type WebApp struct {
	webPath       string
	configuration Configuration
}

func CreateWebApp() WebApp {
	return WebApp{webPath: "/git-stories"}
}

func (this *WebApp) Start() {
	this.loadConfiguration()
	var fileServer = http.FileServer(http.Dir("../frontend/build"))
	var webFilePath = this.webPath + "/static-files"
	http.Handle(webFilePath+"/", http.StripPrefix(webFilePath+"/", fileServer))
	var webApiPath = this.webPath + "/api"
	this.handle(webApiPath+"/repoHistory", this.getRepoHistory)
	this.handle(webApiPath+"/commits", this.commits)
	this.handle(webApiPath+"/fullLog", this.getFullLog)
	this.handle(webApiPath+"/story", this.getStory)
}

func (this *WebApp) loadConfiguration() {
	var configuration, fileError = ioutil.ReadFile("configuration.json")
	if nil != fileError {
		panic(fileError)
	}
	var jsonError = json.Unmarshal(configuration, &this.configuration)
	if nil != jsonError {
		panic(jsonError)
	}
}

func (this *WebApp) handle(path string, function http.HandlerFunc) {
	http.HandleFunc(path, function)
}

func (this *WebApp) getRepoHistory(_ http.ResponseWriter, request *http.Request) {
	var dirPath = request.URL.Query()["dirPath"]
	common.Use(dirPath)
}

func (this *WebApp) commits(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var commits, e = git_client.CreateGitClient(directory).ReadLog(100)
	if e == nil {
		writeJson(responseWriter, commits)
	}
}

func (this *WebApp) getFullLog(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var log, e = git_client.CreateGitClient(directory).ReadDetailedLog(100)
	if nil != e {
		panic(e)
	}
	writeJson(responseWriter, log)
}

func (this *WebApp) getStory(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var lengthLimit = 10
	if len(request.URL.Query()["lengthLimit"]) > 0 {
		var extractedLengthLimit, e = strconv.Atoi(request.URL.Query()["lengthLimit"][0])
		if nil != e {
			panic(e)
		}
		lengthLimit = extractedLengthLimit
	}
	var log, gitError = git_client.CreateGitClient(directory).ReadDetailedLog(lengthLimit)
	if nil != gitError {
		panic(gitError)
	}
	var logBytes, jsonWriteError = json.Marshal(log)
	if nil != jsonWriteError {
		panic(jsonWriteError)
	}
	var workingDirectory, getwdError = os.Getwd()
	if nil != getwdError {
		panic(getwdError)
	}
	var pluginFilePath = workingDirectory + "\\" + this.configuration.Plugin
	var command = exec.Command(pluginFilePath)
	var writer, writerError = command.StdinPipe()
	if nil != writerError {
		panic(writerError)
	}
	var _, writeError = writer.Write(logBytes)
	if nil != writeError {
		panic(writeError)
	}
	var closeError = writer.Close()
	if nil != closeError {
		panic(closeError)
	}
	var output, pluginOutputError = command.CombinedOutput()
	if nil != pluginOutputError {
		panic(pluginOutputError)
	}
	var _, _ = responseWriter.Write(output)
}
