package main

import (
	"fmt"
	"gitwize-be/src/controller"
)

func main() {
	fmt.Println("Hello from gitwize BE")
	r := controller.GetDefaultController()
	r.GET(controller.PingEndPoint, controller.GetPing)
	r.Run()
}
