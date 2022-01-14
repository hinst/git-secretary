package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	log.Fatal(http.ListenAndServe(":3003", nil))
}
