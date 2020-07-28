package controller

import (
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CodeVelocity struct {
	Commits        map[string]string `json:"commits"`
	AdditionLOC    map[string]string `json:"additions"`
	DeletionsLOC   map[string]string `json:"deletions"`
	NewCodeChanges map[string]string `json:"newCodeChanges"`
}

func getCodeChangeVelocity(c *gin.Context) {
	repoID := c.Param("id")

	from, err := strconv.Atoi(c.Query("date_from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	to, err := strconv.Atoi(c.Query("date_to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	commits, additions, deletions, err := db.GetCodeVelocity(repoID, time.Unix(int64(from), 0), time.Unix(int64(to), 0))

	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	velocity := CodeVelocity{
		Commits:      commits,
		AdditionLOC:  additions,
		DeletionsLOC: deletions,
		// Fake newCodeChanges
		NewCodeChanges: commits,
	}

	c.JSON(http.StatusOK, velocity)
	return
}
