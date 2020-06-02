package main

import (
	"fmt"
	"gitwize-be/src/configuration"
	"gitwize-be/src/controller"
)

func main() {
	fmt.Println("Hello from gitwize BE")
	// read configuration file
	configuration.ReadConfiguration()
	// initialize controller
	r := controller.Initialize()
	r.Run()
}
