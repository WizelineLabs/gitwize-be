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

	// rID1 := 61
	// url1 := "git@github.com:go-git/go-git.git"
	// dateRange1 := lambda.GetLastNDayDateRange(30)
	// lambda.UpdateDataForRepo(rID1, url1, "", "", "", dateRange1)

	directory := "/Users/sang.dinh/WorkSpace/GitWize/example-repos/go"
	// url = "git@github.com:golang/go.git"
	// lambda.GitCloneLocal(directory, url)
	rID2 := 62
	dateRange2 := lambda.GetLastNDayDateRange(300)
	lambda.LoadLocalRepo(rID2, directory, "", dateRange2)
}
