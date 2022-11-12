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
	port := configuration.CurConfiguration.Server.Port
	host := configuration.CurConfiguration.Server.Host
	addr := fmt.Sprintf("%s:%s", host, port)

	// initialize controller
	r := controller.Initialize()

	r.Run(addr)
}
