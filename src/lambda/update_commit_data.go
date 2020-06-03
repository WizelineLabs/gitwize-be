package lambda

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"log"
	"os"
	"time"
)

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
	cDtos, fDtos := getDtosFromCommitIter(commitIter, repoID)
	log.Println("\ndto size:", len(cDtos), len(fDtos))
	conn := db.SqlDBConn()
	defer conn.Close()
	executeMulipleBulks(cDtos, conn)
}

func LoadLocalRepo(repoID int, repoPath, branch string, dateRange DateRange) {
	defer timeTrack(time.Now(), "LoadLocalRepo")

	r := getRepoLocal(repoPath)
	commitIter := getCommitIterFromBranch(r, branch, dateRange)
	cDtos, fDtos := getDtosFromCommitIter(commitIter, repoID)
	log.Println("\ndto size:", len(cDtos), len(fDtos))
	conn := db.SqlDBConn()
	executeMulipleBulks(cDtos, conn)
}

func getDtosFromCommitIter(commitIter object.CommitIter, repoID int) (cDtos []commitDto, fDtos []fileStatDTO) {
	defer timeTrack(time.Now(), "getDtosFromCommitIter")

	err := commitIter.ForEach(func(c *object.Commit) error {
		dto := getCommitDTO(c)
		dto.RepositoryID = repoID
		cDtos = append(cDtos, dto)
		fDtos = append(fDtos, getFileStatDTO(c, repoID)...)
		return nil
	})

	if err != nil {
		panic(err.Error())
	}
	return cDtos, fDtos
}
