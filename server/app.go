package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"

	"git-secretary/git_secretary"
)

func main() {
	log.Default().SetPrefix("GS ")
	log.Println("STARTING")
	debug.SetGCPercent(30)

	var workingDirectory = flag.String("wd", "", "Working directory")
	var autoOpenEnabled = flag.Bool("ao", true, "Open URL pointing to the app page automatically on start")
	flag.Parse()
	if workingDirectory != nil && len(*workingDirectory) > 0 {
		log.Println("Go to " + *workingDirectory)
		var e = os.Chdir(*workingDirectory)
		if e != nil {
			panic("Unable to change current directory to " + *workingDirectory)
		}
	}
	var webApp git_secretary.WebApp
	webApp.Create()
	if autoOpenEnabled != nil {
		webApp.Configuration.AutoOpenEnabled = *autoOpenEnabled
	}
	webApp.Start()
}
