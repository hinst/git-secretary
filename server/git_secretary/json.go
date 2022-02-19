package git_secretary

import (
	"encoding/json"
	"github.com/hinst/go-common"
	"net/http"
)

const contentTypeJson = "application/json"

func writeJson(responseWriter http.ResponseWriter, o interface{}) {
	responseWriter.Header().Add(contentTypeHeaderKey, contentTypeJson)
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
