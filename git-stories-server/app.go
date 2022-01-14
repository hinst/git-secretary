package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Now starting...")

	var workingDirectory = flag.String("wd", "", "Working directory")
	flag.Parse()
	if workingDirectory != nil && len(*workingDirectory) > 0 {
		fmt.Println("Go to " + *workingDirectory)
		os.Chdir(*workingDirectory)
	}

	var webApp WebApp
	webApp.Init()
	webApp.Start()
}
