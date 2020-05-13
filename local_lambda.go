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
	// repoID := 10
	// url := "git@github.com:golang/go.git"
	repoID := 2
	url := "git@github.com:go-git/go-git.git"
	dateRange := lambda.GetLastNDayDateRange(3)
	lambda.UpdateDataForRepo(repoID, url, "", "", "", dateRange)
}
