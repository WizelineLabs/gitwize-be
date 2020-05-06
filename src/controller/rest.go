package controller

import "github.com/gin-gonic/gin"

var ginCont *gin.Engine

// GetDefaultController returns Gin Engine Default
func GetDefaultController() *gin.Engine {
	return gin.Default()
}

// GetPing is a handler for /ping endpoint
func GetPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
