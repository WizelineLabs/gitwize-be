package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "gitwize-be/src/auth"
	"gitwize-be/src/controller"
	"os"
)

// AuthMiddleware checks for valid access token
func AuthMiddleware(c *gin.Context) {
	authDisabled := os.Getenv("AUTH_DISABLED") == "true"
	if !authDisabled && !auth.IsAuthorized(nil, c.Request) {
		c.JSON(401, gin.H{
			"message.key": "system.unauthorized",
			"message":     "Unauthorized!",
		})
		c.Abort()
	}

	c.Next()
}

func main() {
	fmt.Println("Hello from gitwize BE")
	r := controller.Initialize()
	// only authorized users can access
	r.Use(AuthMiddleware)
	r.Run()
}
