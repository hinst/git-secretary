package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Now starting...")

	var workingDirectory = flag.String("wd", "", "Working directory")
	flag.Parse()
	if workingDirectory != nil && len(*workingDirectory) > 0 {
		log.Println("Go to " + *workingDirectory)
		os.Chdir(*workingDirectory)
	}

	var webApp WebApp
	webApp.Init()
	webApp.Start()
}
