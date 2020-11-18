package main

import (
	"log"
	"net/http"
)

type WebApp struct {
	webPath string
}

func CreateWebApp() WebApp {
	return WebApp{webPath: "/git-stories"}
}

func (this WebApp) Start() {
	var fileServer = http.FileServer(http.Dir("../frontend/build"))
	var webFilePath = this.webPath + "/static-files"
	http.Handle(webFilePath+"/", http.StripPrefix(webFilePath+"/", fileServer))
	log.Println(webFilePath)
}
