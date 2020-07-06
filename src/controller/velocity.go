package controller

import (
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CodeChangeVelocity struct {
	NetChanges map[string]string `json:"netChanges"`
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

	netChanges, err := db.GetNetChanges(repoID, time.Unix(int64(from), 0), time.Unix(int64(to), 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	velocity := CodeChangeVelocity{
		NetChanges: netChanges,
	}

	c.JSON(http.StatusOK, velocity)
	return
}
