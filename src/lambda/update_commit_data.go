package lambda

import (
	"database/sql"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"log"
	"os"
	"time"
)

func UpdateCommitDataAllRepos() {
	log.Println("start update commit data all repos")
	// find repo id
	conn := db.SqlDBConn()
	rows, _ := conn.Query("SELECT id, url, password, username FROM repository")

	var url, username string
	var id int
	password := sql.NullString{
		String: "",
		Valid:  false,
	}
	if rows == nil {
		log.Printf("[WARN] No repositories found")
		return
	}
	for rows.Next() {
		err := rows.Scan(&id, &url, &password, &username)
		if err != nil {
			log.Panicln(err)
		}
		dateRange := GetLastNDayDateRange(30)
		UpdateDataForRepo(id, url, username, password.String, "", dateRange)
	}
}

// UpdateDataForRepo update data for public/private remote repo using in memory clone
func UpdateDataForRepo(repoID int, repoURL, repoUser, repoPass, branch string, dateRange DateRange) {
	defer timeTrack(time.Now(), "UpdateDataForRepo")
	var r *git.Repository
	if len(repoPass) == 0 {
		r = getPublicRepo(repoURL)
	} else {
		accessToken := cypher.DecryptString(repoPass, os.Getenv("CYPHER_PASS_PHASE"))
		r = getPrivateRepo(repoURL, repoUser, accessToken)
	}
	commitIter := getCommitIterFromBranch(r, branch, dateRange)
	updateCommitData(commitIter, repoID)
}

func updateCommitData(commitIter object.CommitIter, repoID int) {
	defer timeTrack(time.Now(), "updateCommitData")
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

// LoadLocalRepo load data for a local repo already clone on File System
func LoadLocalRepo(repoID int, repoPath, branch string, dateRange DateRange) {
	defer timeTrack(time.Now(), "LoadLocalRepo")
	r := getRepoLocal(repoPath)
	commitIter := getCommitIterFromBranch(r, branch, dateRange)
	updateCommitData(commitIter, repoID)
}
