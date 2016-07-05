package models

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
	"github.com/templecloud/nu/boltdb"
)

//--------------------------------------------------------------------------------------------------

// ResourceName denotes the name used for this tpye of resource in endpoints.
const ResourceName = "functions"

// A Function structure represents an executable funciton definiton.
type Function struct {
	ID             string
	DeployLocation string
	FunctionName   string
	FunctionSrc    string
	Runtime        string
}

// SetID sets the unique id of the specified function
func (function *Function) SetID(ID string) {
	function.ID = ID
}

//--------------------------------------------------------------------------------------------------

// RegisterFunction registers a function
func RegisterFunction(db *bolt.DB, function Function) (Function, error) {

	id := uuid.NewV4()

	function.SetID(id.String())

	functionWithIDBytes, _ := json.Marshal(function)

	boltdb.SetKeyValue(db, []byte(ResourceName), id.Bytes(), functionWithIDBytes)

	persistedFunctionBytes, _ := boltdb.GetValue(db, []byte(ResourceName), id.Bytes())

	persistedFunction := Function{}
	json.Unmarshal(persistedFunctionBytes, &persistedFunction)

	return persistedFunction, nil
}
