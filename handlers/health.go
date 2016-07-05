package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//--------------------------------------------------------------------------------------------------

// Ping writes the string "pong!" back to the client.
func Ping(responseWriter http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	responseWriter.Write([]byte("pong!"))
}

// Echo write the value of the query parameter 'msg' back to the cllent
func Echo(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	toEcho := request.URL.Query().Get("msg")
	responseWriter.Write([]byte(toEcho))
}
