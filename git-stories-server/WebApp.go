package main

import (
	"bufio"
	"encoding/json"
	"git-stories-server/git_client"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hinst/go-common"
	"github.com/pkg/browser"
)

type WebApp struct {
	webPath       string
	configuration Configuration
}

func (me *WebApp) Init() {
	if len(me.webPath) == 0 {
		me.webPath = "/git-stories"
	}
}

func (me *WebApp) Start() {
	me.loadConfiguration()
	var fileServer = http.FileServer(http.Dir("./frontend"))
	var webFilePath = me.webPath + "/static-files"
	http.Handle(webFilePath+"/", http.StripPrefix(webFilePath+"/", fileServer))
	var webApiPath = me.webPath + "/api"
	me.handle(webApiPath+"/repoHistory", me.getRepoHistory)
	me.handle(webApiPath+"/commits", me.commits)
	me.handle(webApiPath+"/fullLog", me.getFullLog)
	me.handle(webApiPath+"/stories", me.getStories)

	var filePicker = FilePicker{WebPath: webApiPath}
	filePicker.Initialize(me.handle)

	me.startListening()
}

func (me *WebApp) startListening() {
	if me.configuration.PortNumber == 0 {
		log.Fatal("Error: Please provide PortNumber in configuration.json")
	}
	var portString = strconv.Itoa(me.configuration.PortNumber)
	var url = "http://localhost:" + portString
	log.Println("Will listen at " + url)
	go func() {
		time.Sleep(1 * time.Second)
		browser.OpenURL(url)
	}()
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}

func (me *WebApp) loadConfiguration() {
	var configuration, fileError = ioutil.ReadFile("configuration.json")
	if nil != fileError {
		panic(fileError)
	}
	var jsonError = json.Unmarshal(configuration, &me.configuration)
	if nil != jsonError {
		panic(jsonError)
	}
}

func (me *WebApp) handle(path string, function http.HandlerFunc) {
	http.HandleFunc(path, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		function(responseWriter, request)
	})
}

func (me *WebApp) getRepoHistory(_ http.ResponseWriter, request *http.Request) {
	var dirPath = request.URL.Query()["dirPath"]
	common.Use(dirPath)
}

func (me *WebApp) commits(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var commits, e = git_client.CreateGitClient(directory).ReadLog(100)
	if e == nil {
		writeJson(responseWriter, commits)
	}
}

func (me *WebApp) getFullLog(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var log, e = git_client.CreateGitClient(directory).ReadDetailedLog(100)
	if nil != e {
		panic(e)
	}
	writeJson(responseWriter, log)
}

func (me *WebApp) getStories(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var lengthLimit = 10
	if len(request.URL.Query()["lengthLimit"]) > 0 {
		var extractedLengthLimit, e = strconv.Atoi(request.URL.Query()["lengthLimit"][0])
		common.AssertError(e)
		lengthLimit = extractedLengthLimit
	}
	var log, gitError = git_client.CreateGitClient(directory).ReadDetailedLog(lengthLimit)
	if gitError != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Unable to open repository at path \"" + directory + "\""))
		return
	}
	var logBytes, jsonWriteError = json.Marshal(log)
	common.AssertError(jsonWriteError)
	var workingDirectory, getwdError = os.Getwd()
	common.AssertError(getwdError)
	var pluginFilePath = workingDirectory + "\\" + me.configuration.Plugin
	var command = exec.Command(pluginFilePath)
	var writer, writerError = command.StdinPipe()
	common.AssertError(writerError)
	var output, pluginOutputError = command.StdoutPipe()
	common.AssertError(pluginOutputError)
	common.AssertError(command.Start())
	var bufferedWriter = bufio.NewWriter(writer)
	var _, bufferedWriteError = bufferedWriter.Write(logBytes)
	common.AssertError(bufferedWriteError, bufferedWriter.Flush())
	var closeError = writer.Close()
	common.AssertError(closeError)
	var outputData, outputError = ioutil.ReadAll(bufio.NewReader(output))
	common.AssertError(outputError)
	responseWriter.Header().Add(contentTypeHeaderKey, contentTypeJson)
	var _, ignoredResponseOutputError = responseWriter.Write(outputData)
	common.Use(ignoredResponseOutputError)
}
