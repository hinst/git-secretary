package main

import (
	"encoding/json"
	"github.com/hinst/go-common"
	"net/http"
)

func writeJson(responseWriter http.ResponseWriter, o interface{}) {
	responseWriter.Header().Add("Content-Type", "application/json")
	var _, e = responseWriter.Write(toJson(o))
	common.Use(e)
}

func toJson(o interface{}) []byte {
	var data, e = json.Marshal(o)
	if e != nil {
		panic(e)
	}
	return data
}
