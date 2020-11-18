package main

import (
	"encoding/json"
	"git-stories/common"
	"git-stories/git_client"
	"net/http"
)

type WebApp struct {
	webPath string
}

func CreateWebApp() WebApp {
	return WebApp{webPath: "/git-stories"}
}

func (this *WebApp) Start() {
	var fileServer = http.FileServer(http.Dir("../frontend/build"))
	var webFilePath = this.webPath + "/static-files"
	http.Handle(webFilePath+"/", http.StripPrefix(webFilePath+"/", fileServer))
	var webApiPath = this.webPath + "/api"
	this.handle(webApiPath+"/repoHistory", this.getRepoHistory)

	this.handle(webApiPath+"/commits", this.commits)
}

func (this *WebApp) handle(path string, function http.HandlerFunc) {
	http.HandleFunc(path, function)
}

func (this *WebApp) getRepoHistory(responseWriter http.ResponseWriter, request *http.Request) {
	var dirPath = request.URL.Query()["dirPath"]
	common.Use(dirPath)
}

func (this *WebApp) commits(responseWriter http.ResponseWriter, request *http.Request) {
	var directory = request.URL.Query()["directory"][0]
	var commits, e = git_client.CreateGitClient(directory).ReadLog()
	if e == nil {
		this.writeJson(responseWriter, commits)
	}
}

func (this *WebApp) writeJson(responseWriter http.ResponseWriter, o interface{}) {
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.Write(toJson(o))
}

func toJson(o interface{}) []byte {
	var data, e = json.Marshal(o)
	if e != nil {
		panic(e)
	}
	return data
}
