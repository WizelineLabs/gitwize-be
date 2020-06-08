package lambda

import (
	"database/sql"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitwize-be/src/configuration"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"gitwize-be/src/utils"
	"log"
	"os"
	"time"
)

//UpdateCommitDataAllRepos update for all repositories from db
func UpdateCommitDataAllRepos() {
	defer utils.TimeTrack(time.Now(), "UpdateCommitDataAllRepos")
	// find repo id
	conn := db.SqlDBConn()
	rows, _ := conn.Query("SELECT id, name, url, password FROM repository")

	var name, url string
	var id int
	password := sql.NullString{
		String: "",
		Valid:  false,
	}
	if rows == nil {
		log.Printf("No repositories found")
		return
	}
	for rows.Next() {
		err := rows.Scan(&id, &name, &url, &password)
		if err != nil {
			log.Panicln(err)
		}
		dateRange := GetLastNDayDateRange(90)
		UpdateDataForRepo(id, url, name, password.String, "", dateRange)
	}
}

// UpdateDataForRepo update data for public/private remote repo using in memory clone
func UpdateDataForRepo(repoID int, repoURL, repoName, repoPass, branch string, dateRange DateRange) {
	defer utils.TimeTrack(time.Now(), "UpdateDataForRepo")

	var r *git.Repository
	var accessToken string
	if mp := os.Getenv("USE_DEFAULT_API_TOKEN"); mp != "" {
		accessToken = os.Getenv("DEFAULT_GITHUB_TOKEN")
	} else {
		accessToken = cypher.DecryptString(repoPass, configuration.CurConfiguration.Cypher.PassPhase)
	}
	r = GetRepo(repoName, repoURL, accessToken)
	commitIter := GetCommitIterFromBranch(r, branch, dateRange)
	updateCommitData(commitIter, repoID)
}

func updateCommitData(commitIter object.CommitIter, repoID int) {
	defer utils.TimeTrack(time.Now(), "updateCommitData")

	conn := db.SqlDBConn()
	defer conn.Close()
	dtos := []commitDto{}
	err := commitIter.ForEach(func(c *object.Commit) error {
		if len(dtos) == batchSize {
			executeBulkStatement(dtos, conn)
			dtos = []commitDto{}
		} else {
			dto := getCommitDTO(c)
			dto.RepositoryID = repoID
			dtos = append(dtos, dto)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}
	if len(dtos) > 0 {
		executeBulkStatement(dtos, conn)
	}
}
