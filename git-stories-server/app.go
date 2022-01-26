package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	fmt.Println("Now starting...")
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
	var webApp WebApp
	webApp.Create()
	if autoOpenEnabled != nil {
		log.Println(*autoOpenEnabled)
		webApp.configuration.AutoOpenEnabled = *autoOpenEnabled
	}
	webApp.Start()
}
