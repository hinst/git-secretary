package git_secretary

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	"github.com/pkg/browser"
	bolt "go.etcd.io/bbolt"
)

type WebApp struct {
	Configuration Configuration
	Storage       Storage
	storage       *bolt.DB
	tasks         *WebTaskManager
}

const FILE_PERMISSION_OWNER_READ_WRITE = 0600

func (me *WebApp) Create() {
	me.Configuration.SetDefault()
	me.loadConfiguration()

	me.Storage.Create()
	var dbOptions = *bolt.DefaultOptions
	dbOptions.Timeout = 1 * time.Second
	dbOptions.ReadOnly = false
	var storage, e = bolt.Open("./storage.bolt", FILE_PERMISSION_OWNER_READ_WRITE, &dbOptions)
	common.AssertError(common.CreateExceptionIf("Unable to open storage file", e))
	me.storage = storage

	me.tasks = (&WebTaskManager{}).Create()
}

func (me *WebApp) GetWebFilePath() string {
	return me.Configuration.WebPath + "/static-files"
}

func (me *WebApp) Start() {
	var fileServer = http.FileServer(http.Dir("./frontend"))
	var webFilePath = me.GetWebFilePath()
	http.Handle(webFilePath+"/", http.StripPrefix(webFilePath+"/", fileServer))
	var webApiPath = me.Configuration.WebPath + "/api"
	me.handle(webApiPath+"/repoHistory", me.getRepoHistory)
	me.handle(webApiPath+"/commits", me.commits)
	me.handle(webApiPath+"/fullLog", me.getFullLog)
	me.handle(webApiPath+"/stories", me.getStoriesAsync)
	me.handle(webApiPath+"/task", me.getTask)

	var filePicker = FilePicker{WebPath: webApiPath}
	filePicker.Initialize(me.handle)

	me.startListening()
}

func (me *WebApp) startListening() {
	if me.Configuration.PortNumber == 0 {
		log.Fatal("Error: Please provide PortNumber in configuration.json")
	}
	var portString = strconv.Itoa(me.Configuration.PortNumber)
	var url = "http://localhost:" + portString + me.GetWebFilePath()
	log.Println("Will listen at " + url)
	if me.Configuration.AutoOpenEnabled {
		go func() {
			time.Sleep(1 * time.Second)
			browser.OpenURL(url)
		}()
	}
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}

func (me *WebApp) loadConfiguration() {
	var configuration, fileError = ioutil.ReadFile("configuration.json")
	if nil != fileError {
		panic(fileError)
	}
	var jsonError = json.Unmarshal(configuration, &me.Configuration)
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

func (me *WebApp) requireDirectoryArgument(responseWriter http.ResponseWriter, request *http.Request) (directory string) {
	directory = request.URL.Query().Get("directory")
	if len(directory) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Query argument \"directory\" is required"))
	}
	return
}

func (me *WebApp) getStoriesAsync(responseWriter http.ResponseWriter, request *http.Request) {
	var requestObject ReadStoriesRequest
	requestObject.Directory = me.requireDirectoryArgument(responseWriter, request)
	if len(requestObject.Directory) == 0 {
		return
	}
	requestObject.TimeZone = request.URL.Query().Get("timeZone")
	var lengthLimitText = request.URL.Query().Get("lengthLimit")
	if len(lengthLimitText) > 0 {
		var parsedLengthLimit, e = strconv.Atoi(lengthLimitText)
		if nil != e {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write([]byte("Unable to parse query argument lengthLimit as integer: \"" + lengthLimitText + "\""))
			return
		}
		requestObject.LengthLimit = parsedLengthLimit
	}
	var taskId = me.tasks.Add(&WebTask{})
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(strconv.FormatUint(uint64(taskId), 10)))
	go me.readStories(taskId, requestObject)
}

func (me *WebApp) readStories(taskId uint, request ReadStoriesRequest) {
	var gitClient = (&CachedGitClient{}).Create(me.storage, request.Directory)
	gitClient.SetProgressReceiver(func(total int, done int) {
		me.tasks.Update(taskId, func(task *WebTask) {
			task.Total = total
			task.Done = done
		})
	})
	var rows, gitError = gitClient.ReadDetailedLog(0)
	if nil != gitError {
		me.tasks.Update(taskId, func(task *WebTask) {
			task.Error = gitError.Error()
		})
		return
	}
	var workingDirectory, getwdError = os.Getwd()
	common.AssertError(getwdError)
	var pluginFilePath = workingDirectory + "/" + me.Configuration.Plugin
	var pluginRunner = PluginRunner{PluginFilePath: pluginFilePath}
	var storyEntries, pluginError = pluginRunner.Run(git_stories_api.StoriesRequest{
		LogEntries: rows,
		TimeZone:   request.TimeZone,
	})
	if nil != pluginError {
		me.tasks.Update(taskId, func(task *WebTask) {
			task.Error = pluginError.Error()
		})
		return
	}
	me.tasks.Update(taskId, func(task *WebTask) {
		task.StoryEntries = storyEntries
		if nil == task.StoryEntries {
			// Avoid nil value because nil means that the task is not finished yet
			task.StoryEntries = make([]git_stories_api.StoryEntryChangeset, 0)
		}
		if request.LengthLimit > 0 && len(task.StoryEntries) > request.LengthLimit {
			task.StoryEntries = task.StoryEntries[0:request.LengthLimit]
		}
	})
}

func (me *WebApp) getTask(responseWriter http.ResponseWriter, request *http.Request) {
	var idString = request.URL.Query().Get("id")
	if len(idString) == 0 {
		responseWriter.WriteHeader(http.StatusBadGateway)
		responseWriter.Write([]byte("Query parameter is required: id"))
		return
	}
	var id64, idParseError = strconv.ParseUint(idString, 10, int(common.SizeOfUint))
	if nil != idParseError {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("Query parameter must be an unsigned integer: id; got: " +
			idString + "\n" + idParseError.Error()))
		return
	}
	var id = uint(id64)
	var task = me.tasks.Get(id)
	if nil == task {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte("Task not found: id=" + strconv.FormatUint(id64, 10)))
		return
	}
	var taskBytes, jsonError = json.Marshal(task)
	if jsonError != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Unable to encode json\n" + jsonError.Error()))
	}
	responseWriter.Header().Add(contentTypeHeaderKey, contentTypeJson)
	responseWriter.Write(taskBytes)
}