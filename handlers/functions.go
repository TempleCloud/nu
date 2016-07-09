package handlers

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/templecloud/nu/models"
)

//--------------------------------------------------------------------------------------------------

// ListFunctions list all functions
func ListFunctions(db *bolt.DB) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			allFunctions, _ := models.ListFunctions(db)
			c.JSON(http.StatusOK, allFunctions)
		})
}

// GetFunction gets the specified function
func GetFunction(db *bolt.DB) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			function, _ := models.GetFunction(db, c.Param("functionId"))
			c.JSON(http.StatusOK, function)
		})
}

// RegisterFunction registers a function
func RegisterFunction(db *bolt.DB) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			// Extract input function to register
			var function models.Function
			c.BindJSON(&function)
			// Return registered function
			registeredFunction, _ := models.RegisterFunction(db, function)
			c.JSON(http.StatusOK, registeredFunction)
		})
}

// UpdateFunction updates a registered function
func UpdateFunction(db *bolt.DB) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			// Extract input function to register
			var function models.Function
			c.BindJSON(&function)
			// Return registered function
			registeredFunction, err := models.UpdateFunction(db, c.Param("functionId"), function)
			if err == nil {
				c.JSON(http.StatusOK, registeredFunction)
			} else {
				c.AbortWithError(http.StatusBadRequest, err)
			}
		})
}

// DeleteFunction deletes the specified function
func DeleteFunction(db *bolt.DB) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			err := models.DeleteFunction(db, c.Param("functionId"))
			if err != nil {
				c.Status(http.StatusOK)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		})
}
