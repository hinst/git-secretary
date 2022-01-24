package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
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

func (me *WebApp) GetWebFilePath() string {
	return me.webPath + "/static-files"
}

func (me *WebApp) Start() {
	me.loadConfiguration()
	var fileServer = http.FileServer(http.Dir("./frontend"))
	var webFilePath = me.GetWebFilePath()
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
	var url = "http://localhost:" + portString + me.GetWebFilePath()
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
	var commits, e = CreateGitClient(directory).ReadLog(100)
	if e == nil {
		writeJson(responseWriter, commits)
	}
}

func (me *WebApp) getFullLog(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var log, e = CreateGitClient(directory).ReadDetailedLog(100)
	if nil != e {
		panic(e)
	}
	writeJson(responseWriter, log)
}

func (me *WebApp) getStories(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	if len(directory) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Query argument \"directory\" is required"))
		return
	}
	var lengthLimit = math.MaxInt32
	if len(request.URL.Query()["lengthLimit"]) > 0 {
		var extractedLengthLimit, e = strconv.Atoi(request.URL.Query()["lengthLimit"][0])
		common.AssertError(e)
		lengthLimit = extractedLengthLimit
	}
	var rows, gitError = CreateGitClient(directory).
		SetDebugLogEnabled(true).
		ReadDetailedLog(lengthLimit)
	if gitError != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Unable to open repository at path \"" + directory + "\""))
		return
	}
	var workingDirectory, getwdError = os.Getwd()
	common.AssertError(getwdError)
	var pluginFilePath = workingDirectory + "/" + me.configuration.Plugin
	var pluginRunner = PluginRunner{PluginFilePath: pluginFilePath}
	var storyEntries = pluginRunner.Run(rows)
	var outputBytes, jsonError = json.Marshal(storyEntries)
	common.AssertWrapped(jsonError, "Unable to encode json: storyEntries")
	responseWriter.Header().Add(contentTypeHeaderKey, contentTypeJson)
	var _, responseWriteError = responseWriter.Write(outputBytes)
	common.Use(responseWriteError)
}
