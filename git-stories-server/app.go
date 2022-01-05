package main

import (
	"log"
	"net/http"
)

func main() {
	println("Now starting...")
	var webApp = CreateWebApp()
	webApp.Start()
	log.Fatal(http.ListenAndServe(":3003", nil))
}
