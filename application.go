package main

import (
	"fmt"
	"gitwize-be/src/controller"
)

func main() {
	fmt.Println("Hello from gitwize BE")
	r := controller.Initialize()
	// only authorized users can access
	r.Run()
}
