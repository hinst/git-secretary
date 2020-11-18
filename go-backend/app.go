package main

import "net/http"
import "log"

func main() {
	println("Now starting...")
	var webApp = CreateWebApp()
	webApp.Start()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
