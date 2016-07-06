package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/templecloud/nu/models"
)

//--------------------------------------------------------------------------------------------------

// ListFunctions list all functions
func ListFunctions(db *bolt.DB) httprouter.Handle {
	return httprouter.Handle(
		func(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
			allFunctions, _ := models.ListFunctions(db)
			allFunctionsBytes, _ := json.Marshal(allFunctions)
			responseWriter.Write(allFunctionsBytes)
		})
}

// RegisterFunction registers a function
func RegisterFunction(db *bolt.DB) httprouter.Handle {
	return httprouter.Handle(
		func(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {

			functionBytes, _ := ioutil.ReadAll(request.Body)
			function := models.Function{}
			json.Unmarshal(functionBytes, &function)

			registeredFunction, _ := models.RegisterFunction(db, function)

			registeredFunctionBytes, _ := json.Marshal(registeredFunction)
			responseWriter.Write(registeredFunctionBytes)
		})
}
