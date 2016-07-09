package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//--------------------------------------------------------------------------------------------------

// Ping writes the string "pong!" back to the client.
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong!")
}

// Echo write the value of the query parameter 'msg' back to the cllent
func Echo(c *gin.Context) {
	msg := c.Query("msg")
	c.String(http.StatusOK, msg)
}
