/*
this is local run example
should be removed after lambda deployment completed
to set up:
	- setup environment vars for DB connection
	- change main2() => main()
	- go run local_lambda.go
*/

package main

import (
	"gitwize-be/src/lambda"
)

func main2() {
	lambda.UpdateCommitDataAllRepos()
}
