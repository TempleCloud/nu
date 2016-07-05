package functions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

//--------------------------------------------------------------------------------------------------

// RegisterFunction registers a function
func RegisterFunction(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	// Extract JSON body
	// Generate an id
	// Add to database

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	bodyString := string(bodyBytes)
	log.Printf("function str: %v", bodyString)

	function := models.Function{}
	json.Unmarshal([]byte(bodyString), &function)
	log.Printf("function json: %v", function)

	id := uuid.NewV4()
	function.SetID(id.String())
	log.Printf("updated function json: %v", function)

	updatedFunctionBytes, _ := json.Marshal(function)
	log.Printf("updated function str: %v", string(updatedFunctionBytes))

	models.SetKeyValue(db, []byte(models.ResourceName), id.Bytes(), updatedFunctionBytes)

	saved, _ := models.GetValue(db, []byte(models.ResourceName), id.Bytes())
	log.Printf("saved: %v", saved)

	models.

		// responseWriter.Write([]byte("Not implemented"))
		responseWriter.Write([]byte(saved))
}
