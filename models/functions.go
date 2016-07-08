package models

import (
	"encoding/json"
	"fmt"

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

// ListFunctions registers a function
func ListFunctions(db *bolt.DB) ([]Function, error) {

	var result = make([]Function, 0)

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte(ResourceName))

		bucket.ForEach(func(key, value []byte) error {
			fmt.Printf("key=%s, value=%s\n", key, value)

			persistedFunction := Function{}
			json.Unmarshal(value, &persistedFunction)

			result = append(result, persistedFunction)
			fmt.Printf("result: %v\n", result)

			return nil
		})
		return nil
	})

	return result, nil
}

// GetFunction gets a function definition
func GetFunction(db *bolt.DB, id string) (Function, error) {

	var function Function

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ResourceName))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(ResourceName))
		}

		functionBytes := bucket.Get([]byte(id))
		persistedFunction := Function{}
		json.Unmarshal(functionBytes, &persistedFunction)

		function = persistedFunction

		return nil
	})

	return function, nil
}

// RegisterFunction registers a function
func RegisterFunction(db *bolt.DB, function Function) (Function, error) {

	id := uuid.NewV4().String()

	function.SetID(id)

	functionWithIDBytes, _ := json.Marshal(function)

	boltdb.SetKeyValue(db, []byte(ResourceName), []byte(id), functionWithIDBytes)

	persistedFunctionBytes, _ := boltdb.GetValue(db, []byte(ResourceName), []byte(id))

	persistedFunction := Function{}
	json.Unmarshal(persistedFunctionBytes, &persistedFunction)

	return persistedFunction, nil
}
