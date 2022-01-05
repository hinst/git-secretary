package main

import (
	"log"
	"net/http"
)

func main() {
	println("Now starting...")
	var webApp WebApp
	webApp.Init()
	webApp.Start()
	log.Fatal(http.ListenAndServe(":3003", nil))
}
